package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationServiceServer) CreateReservation(ctx context.Context, req *guestProto.CreateReservationRequest) (*guestProto.CreateReservationResponse, error) {
	res, err := s.reservationService.CreateReservation(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
