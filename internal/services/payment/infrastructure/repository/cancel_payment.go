package repository

import (
	"context"
	"net/http"
	"strings"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	invoiceDomain "github.com/goplaceapp/goplace-guest/internal/services/payment/domain"
	"github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	"google.golang.org/grpc/status"
)

func (r *PaymentRepository) CancelPayment(ctx context.Context, req *guestProto.CancelPaymentRequest) (*guestProto.CancelPaymentResponse, error) {
	var (
		reservation *domain.Reservation
		invoice     *invoiceDomain.Invoice
	)

	if err := r.GetTenantDBConnection(ctx).Table("reservations").Joins("JOIN payment_requests ON payment_requests.reservation_id = reservations.id").Joins("JOIN invoices ON invoices.payment_request_id = payment_requests.id").Where("invoices.id = ?", req.GetInvoiceId()).First(&reservation).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if err := r.GetTenantDBConnection(ctx).First(&invoice, "id = ?", req.GetInvoiceId()).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	integration, err := r.commonRepository.GetIntegrationBySystemType(ctx, "Payment Gateway", reservation.BranchID)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, "Integration not found")
	}

	switch strings.ToLower(integration.SystemName) {
	case "tap":
		if err := r.CancelTapInvoice(invoice.InvoiceID, integration); err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	default:
		return nil, status.Error(http.StatusInternalServerError, "Integration not found")
	}

	if err := r.GetTenantDBConnection(ctx).
		Table("invoices").
		Where("invoice_id = ?", req.GetInvoiceId()).
		Update("status", "canceled").Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &guestProto.CancelPaymentResponse{
		Code:    http.StatusOK,
		Message: "Invoice canceled",
	}, nil
}
