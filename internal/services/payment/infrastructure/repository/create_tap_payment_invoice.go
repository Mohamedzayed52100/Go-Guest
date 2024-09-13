package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-common/pkg/utils"
	"github.com/goplaceapp/goplace-guest/internal/services/payment/domain"
	integrationDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation/domain"
	"github.com/nyaruka/phonenumbers"
	"google.golang.org/grpc/status"
)

func (r *PaymentRepository) CreateTapPaymentInvoice(ctx context.Context, payment *domain.PaymentRequest, integration *integrationDomain.Integration) (*domain.PaymentRequest, map[string]interface{}, error) {
	var (
		tapIntegrationCredentials map[string]interface{}
		parseCredentials          = []byte(integration.Credentials)
	)

	if err := json.Unmarshal(parseCredentials, &tapIntegrationCredentials); err != nil {
		return nil, nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	parsedGuestPhoneNumber, err := phonenumbers.Parse(payment.Guest.PhoneNumber, "")
	if err != nil {
		return nil, nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	regionNumber := phonenumbers.GetRegionCodeForNumber(parsedGuestPhoneNumber)
	countryCode := phonenumbers.GetCountryCodeForRegion(regionNumber)

	var clientId string
	if err := r.GetSharedDB().
		Table("tenant_credentials").
		Where("upper(name) = 'TAP PAYMENT' and enabled = true").
		Select("client_id").
		Scan(&clientId).Error; err != nil {
		return nil, nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	payloadItems := []map[string]interface{}{}
	for _, item := range payment.Items {
		payloadItems = append(payloadItems, map[string]interface{}{
			"name":     item.Name,
			"quantity": item.Quantity,
			"amount":   item.Price,
		})
	}

	var webhookUrl string
	switch os.Getenv("ENVIRONMENT") {
	case meta.ProdEnvironment:
		webhookUrl = "https://api.goplace.io/api/v1/webhooks/payment/" + clientId
	case meta.StagingEnvironment:
		webhookUrl = "https://api-staging.goplace.io/api/v1/webhooks/payment/" + clientId
	default:
		webhookUrl = "https://api-dev.goplace.io/api/v1/webhooks/payment/" + clientId
	}

	total, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", payment.Invoice.SubTotal), 32)

	payload := map[string]interface{}{
		"due":    time.Now().UTC().AddDate(0, 0, 5).UnixMilli(),
		"expiry": time.Now().UTC().AddDate(0, 0, 5).UnixMilli(),
		"mode":   "INVOICE",
		"notifications": map[string]interface{}{
			"channels": []string{},
		},
		"customer": map[string]interface{}{
			"first_name": payment.Guest.FirstName,
			"last_name":  payment.Guest.LastName,
			"email":      payment.Guest.Email,
			"phone": map[string]interface{}{
				"country_code": "+" + strconv.Itoa(countryCode),
				"number":       payment.Guest.PhoneNumber[len(strconv.Itoa(countryCode)):],
			},
		},
		"order": map[string]interface{}{
			"currency": "SAR",
			"amount":   float32(total),
			"items":    payloadItems,
		},
		"post": map[string]interface{}{
			"url": webhookUrl,
		},
	}

	jsonString, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	request, _ := http.NewRequest("POST", integration.BaseURL+"/invoices/", bytes.NewBuffer(jsonString))
	request.Header.Set("Authorization", "Bearer "+utils.Decrypt(tapIntegrationCredentials["secret_key"].(string), os.Getenv("AES_ENCRYPTION_KEY")))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Default().Error(err)
		return nil, nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Default().Error(err)
		return nil, nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	var jsonInvoiceRes map[string]interface{}
	err = json.Unmarshal(body, &jsonInvoiceRes)
	if err != nil {
		logger.Default().Error(err)
		return nil, nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if jsonInvoiceRes["errors"] != nil {
		logger.Default().Error(jsonInvoiceRes["errors"].([]interface{})[0].(map[string]interface{})["description"].(string))
		return nil, nil, status.Error(http.StatusInternalServerError, jsonInvoiceRes["errors"].([]interface{})[0].(map[string]interface{})["description"].(string))
	}

	payment.Invoice = &domain.Invoice{
		InvoiceID:        jsonInvoiceRes["id"].(string),
		PaymentRequestID: payment.ID,
		Status:           "unpaid",
		CustomerID:       jsonInvoiceRes["customer"].(map[string]interface{})["id"].(string),
		Currency:         jsonInvoiceRes["order"].(map[string]interface{})["currency"].(string),
		SubTotal:         payment.Invoice.SubTotal,
	}

	if err := r.GetTenantDBConnection(ctx).Create(&payment.Invoice).Error; err != nil {
		logger.Default().Error(err)
		return nil, nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return payment, jsonInvoiceRes, nil
}
