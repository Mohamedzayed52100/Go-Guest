package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationServiceServer) RequestReservationWebhook(ctx context.Context, req *guestProto.RequestReservationWebhookRequest) (*guestProto.RequestReservationWebhookResponse, error) {
	res, err := s.reservationService.RequestReservationWebhook(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil

}
