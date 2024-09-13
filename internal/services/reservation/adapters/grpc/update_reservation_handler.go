package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationServiceServer) UpdateReservation(ctx context.Context, req *guestProto.UpdateReservationRequest) (*guestProto.UpdateReservationResponse, error) {
	res, err := s.reservationService.UpdateReservation(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *ReservationServiceServer) UpdateReservationFromWebhook(ctx context.Context, req *guestProto.UpdateReservationFromWebhookRequest) (*guestProto.UpdateReservationFromWebhookResponse, error) {
	res, err := s.reservationService.UpdateReservationFromWebhook(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
