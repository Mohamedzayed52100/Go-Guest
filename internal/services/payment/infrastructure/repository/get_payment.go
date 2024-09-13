package repository

import (
	"context"
	"net/http"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/payment/adapters/converters"
	"github.com/goplaceapp/goplace-guest/internal/services/payment/domain"
	guestDomain "github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
	"google.golang.org/grpc/status"
)

func (r *PaymentRepository) GetPaymentByID(ctx context.Context, req *guestProto.GetPaymentByIDRequest) (*guestProto.GetPaymentByIDResponse, error) {
	var (
		payment *domain.PaymentRequest
		err     error
	)

	if err = r.GetTenantDBConnection(ctx).Where("id =? AND reservation_id = ?", req.GetId(), req.GetReservationId()).First(&payment).Error; err != nil {
		return nil, status.Error(http.StatusNotFound, "Payment not found")
	}

	payment, err = r.GetAllPaymentData(ctx, payment)
	if err != nil {
		return nil, err
	}

	return &guestProto.GetPaymentByIDResponse{
		Result: converters.BuildPaymentResponse(payment),
	}, nil
}

func (r *PaymentRepository) GetAllReservationPayments(ctx context.Context, req *guestProto.GetAllReservationPaymentsRequest) (*guestProto.GetAllReservationPaymentsResponse, error) {
	var (
		payments = []*domain.PaymentRequest{}
		err      error
	)

	if err = r.GetTenantDBConnection(ctx).Where("reservation_id = ?", req.GetId()).Order("created_at desc").Find(&payments).Error; err != nil {
		return nil, status.Error(http.StatusNotFound, "Payment not found")
	}

	for i := range payments {
		payments[i], err = r.GetAllPaymentData(ctx, payments[i])
		if err != nil {
			return nil, err
		}
	}

	return &guestProto.GetAllReservationPaymentsResponse{
		Result: converters.BuildAllPaymentsResponse(payments),
	}, nil
}

func (r *PaymentRepository) GetAllPaymentData(ctx context.Context, payment *domain.PaymentRequest) (*domain.PaymentRequest, error) {
	var guest *guestDomain.Guest

	if err := r.GetTenantDBConnection(ctx).
		Table("guests").
		Joins("JOIN reservations ON reservations.guest_id = guests.id").
		Where("reservations.id = ?", payment.ReservationID).
		Select("first_name, last_name, phone_number").
		Scan(&guest).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, "Failed to get guest data")
	}

	payment.Guest = &domain.PaymentGuest{
		FirstName:   guest.FirstName,
		LastName:    guest.LastName,
		PhoneNumber: guest.PhoneNumber,
		Address:     guest.Address,
	}

	branch, err := r.userClient.Client.UserService.Repository.GetBranchByID(ctx, r.userClient.Client.UserService.Repository.GetCurrentBranchId(ctx))
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	payment.Branch = &domain.PaymentBranch{
		Name:          branch.Name,
		Address:       branch.Address,
		VatPercent:    branch.VatPercent,
		ServiceCharge: branch.ServiceCharge,
		CrNumber:      branch.CrNumber,
		VatRegNumber:  branch.VatRegNumber,
	}

	r.GetTenantDBConnection(ctx).
		Unscoped().
		Table("payment_item_assignments").
		Joins("JOIN restaurant_items ON restaurant_items.id = payment_item_assignments.item_id").
		Select("restaurant_items.id, restaurant_items.name, payment_item_assignments.price, payment_item_assignments.quantity").
		Where("payment_id = ?", payment.ID).
		Find(&payment.Items)

	if err := r.GetTenantDBConnection(ctx).
		First(&payment.Invoice, "payment_request_id = ?", payment.ID).
		Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, "Failed to get invoice data")
	}

	return payment, nil
}
