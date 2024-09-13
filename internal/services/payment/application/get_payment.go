package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *PaymentService) GetPaymentByID(ctx context.Context, req *guestProto.GetPaymentByIDRequest) (*guestProto.GetPaymentByIDResponse, error) {
	res, err := s.Repository.GetPaymentByID(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *PaymentService) GetAllReservationPayments(ctx context.Context, req *guestProto.GetAllReservationPaymentsRequest) (*guestProto.GetAllReservationPaymentsResponse, error) {
	res, err := s.Repository.GetAllReservationPayments(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
