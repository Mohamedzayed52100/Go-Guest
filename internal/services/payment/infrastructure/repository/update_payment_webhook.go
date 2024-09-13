package repository

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-common/pkg/utils"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/payment/domain"
	tagDomain "github.com/goplaceapp/goplace-settings/pkg/reservationtagservice/domain"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (r *PaymentRepository) UpdatePaymentFromWebhook(ctx context.Context, req *guestProto.UpdatePaymentFromWebhookRequest) (*guestProto.UpdatePaymentFromWebhookResponse, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for k, v := range md {
			if strings.EqualFold(k, "clientId") {
				var tenantDbName string
				if err := r.GetSharedDB().
					Table("tenant_credentials").
					Joins("join tenants on tenants.id = tenant_credentials.tenant_id").
					Where("client_id = ? and enabled = true", v).
					Select("db_name").
					Scan(&tenantDbName).Error; err != nil {
					return nil, status.Error(http.StatusInternalServerError, "Invalid client id")
				}

				ctx = context.WithValue(ctx, meta.TenantDBNameContextKey.String(), tenantDbName)
			}
		}
	} else {
		return nil, status.Error(http.StatusInternalServerError, "Invalid client id")
	}

	var branchId int32
	if err := r.GetTenantDBConnection(ctx).
		Table("invoices").
		Joins("JOIN payment_requests ON payment_requests.id = invoices.payment_request_id").
		Joins("JOIN reservations ON reservations.id = payment_requests.reservation_id").
		Where("invoices.invoice_id = ?", req.GetInvoiceId()).
		Select("reservations.branch_id").
		Scan(&branchId).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	integration, err := r.commonRepository.GetIntegrationBySystemType(ctx, "Payment Gateway", branchId)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, "Integration not found")
	}

	var (
		tapIntegrationCredentials map[string]interface{}
		parseCredentials          = []byte(integration.Credentials)
	)

	if err := json.Unmarshal(parseCredentials, &tapIntegrationCredentials); err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	request, _ := http.NewRequest("GET", "https://api.tap.company/v2/"+"/invoices/"+req.GetInvoiceId(), nil)
	request.Header.Set("Authorization", "Bearer "+utils.Decrypt(tapIntegrationCredentials["secret_key"].(string), os.Getenv("AES_ENCRYPTION_KEY")))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	res, _ := http.DefaultClient.Do(request)
	body, _ := io.ReadAll(res.Body)

	var jsonInvoiceRes map[string]interface{}
	err = json.Unmarshal(body, &jsonInvoiceRes)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	defer res.Body.Close()

	if strings.EqualFold(req.GetStatus(), jsonInvoiceRes["status"].(string)) {
		if err := r.GetTenantDBConnection(ctx).Where("invoice_id = ?", req.GetInvoiceId()).Updates(&domain.Invoice{
			Status:         req.GetStatus(),
			LastFourDigits: req.GetCard().GetFourDigits(),
			CardType:       req.GetCard().GetBrand(),
		}).Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	} else {
		return nil, status.Error(http.StatusInternalServerError, "Wrong status provided")
	}

	var paidTag *tagDomain.ReservationTag
	r.GetTenantDBConnection(ctx).
		First(&paidTag, "UPPER(name) = UPPER('Paid By Tap') AND branch_id = ?", branchId)

	if paidTag != nil {
		var reservationId int32
		if err := r.GetTenantDBConnection(ctx).Table("reservations").Joins("JOIN payment_requests ON payment_requests.reservation_id = reservations.id").Joins("JOIN invoices ON invoices.payment_request_id = payment_requests.id").Where("invoices.invoice_id = ?", req.GetInvoiceId()).Select("reservations.id").Scan(&reservationId).Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		if err := r.GetTenantDBConnection(ctx).Create(&tagDomain.ReservationTagsAssignment{
			ReservationID: reservationId,
			TagID:         paidTag.ID,
		}).Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	return &guestProto.UpdatePaymentFromWebhookResponse{
		Code:    http.StatusOK,
		Message: "Payment updated successfully",
	}, nil
}
