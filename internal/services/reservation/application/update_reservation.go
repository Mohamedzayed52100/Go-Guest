package application

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationService) UpdateReservation(ctx context.Context, req *guestProto.UpdateReservationRequest) (*guestProto.UpdateReservationResponse, error) {
	res, err := s.Repository.UpdateReservation(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *ReservationService) UpdateReservationFromWebhook(ctx context.Context, req *guestProto.UpdateReservationFromWebhookRequest) (*guestProto.UpdateReservationFromWebhookResponse, error) {
	res, err := s.Repository.UpdateReservationFromWebhook(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
