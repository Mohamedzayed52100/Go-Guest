package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationWaitListServiceServer) CreateDirectInReservation(ctx context.Context, req *guestProto.CreateWaitingReservationRequest) (*guestProto.CreateReservationResponse, error) {
	res, err := s.reservationWaitListService.CreateDirectInReservation(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
