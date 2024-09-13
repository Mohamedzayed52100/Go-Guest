package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationService) RequestReservationWebhook(ctx context.Context, req *guestProto.RequestReservationWebhookRequest) (*guestProto.RequestReservationWebhookResponse, error) {
	res, err := s.Repository.RequestReservationWebhook(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
