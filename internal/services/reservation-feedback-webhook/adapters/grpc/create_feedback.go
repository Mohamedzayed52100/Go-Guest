package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationFeedbackWebhookServiceServer) CreateReservationFeedbackFromWebhook(ctx context.Context, req *guestProto.CreateReservationFeedbackFromWebhookRequest) (*guestProto.CreateReservationFeedbackFromWebhookResponse, error) {
	res, err := s.reservationFeedbackWebhookService.CreateReservationFeedbackFromWebhook(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
