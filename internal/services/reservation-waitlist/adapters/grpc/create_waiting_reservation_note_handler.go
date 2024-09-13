package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationWaitListServiceServer) CreateWaitingReservationNote(ctx context.Context, req *guestProto.CreateWaitingReservationNoteRequest) (*guestProto.CreateWaitingReservationNoteResponse, error) {
	res, err := s.reservationWaitListService.CreateWaitingReservationNote(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
