package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *PaymentService) CancelPayment(ctx context.Context, req *guestProto.CancelPaymentRequest) (*guestProto.CancelPaymentResponse, error) {
	res, err := s.Repository.CancelPayment(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
