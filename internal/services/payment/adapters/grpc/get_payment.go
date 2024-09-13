package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *PaymentServiceServer) GetPaymentByID(ctx context.Context, req *guestProto.GetPaymentByIDRequest) (*guestProto.GetPaymentByIDResponse, error) {
	res, err := s.paymentService.GetPaymentByID(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *PaymentServiceServer) GetAllReservationPayments(ctx context.Context, req *guestProto.GetAllReservationPaymentsRequest) (*guestProto.GetAllReservationPaymentsResponse, error) {
	res, err := s.paymentService.GetAllReservationPayments(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
