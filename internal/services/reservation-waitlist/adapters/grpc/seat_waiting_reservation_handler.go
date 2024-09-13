package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationWaitListServiceServer) SeatWaitingReservation(ctx context.Context, req *guestProto.SeatWaitingReservationRequest) (*guestProto.SeatWaitingReservationResponse, error) {
	res, err := s.reservationWaitListService.SeatWaitingReservation(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
