package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *PaymentService) UpdatePaymentFromWebhook(ctx context.Context, req *guestProto.UpdatePaymentFromWebhookRequest) (*guestProto.UpdatePaymentFromWebhookResponse, error) {
	res, err := s.Repository.UpdatePaymentFromWebhook(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
