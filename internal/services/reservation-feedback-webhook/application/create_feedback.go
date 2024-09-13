package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationFeedbackWebhookService) CreateReservationFeedbackFromWebhook(ctx context.Context, req *guestProto.CreateReservationFeedbackFromWebhookRequest) (*guestProto.CreateReservationFeedbackFromWebhookResponse, error) {
	res, err := s.Repository.CreateReservationFeedbackFromWebhook(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
