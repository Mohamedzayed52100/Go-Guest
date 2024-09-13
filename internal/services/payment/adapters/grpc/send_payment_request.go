package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *PaymentServiceServer) SendPaymentRequest(ctx context.Context, req *guestProto.PaymentRequest) (*guestProto.PaymentResponse, error) {
	res, err := s.paymentService.SendPaymentRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
