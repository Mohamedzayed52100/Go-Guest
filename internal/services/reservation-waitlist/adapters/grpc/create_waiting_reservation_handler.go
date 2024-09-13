package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationWaitListServiceServer) CreateWaitingReservation(ctx context.Context, req *guestProto.CreateWaitingReservationRequest) (*guestProto.CreateWaitingReservationResponse, error) {
	res, err := s.reservationWaitListService.CreateWaitingReservation(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
