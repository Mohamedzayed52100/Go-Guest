package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *PaymentService) SendPaymentRequest(ctx context.Context, req *guestProto.PaymentRequest) (*guestProto.PaymentResponse, error) {
	res, err := s.Repository.SendPaymentRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
