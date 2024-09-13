package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationWaitListServiceServer) UpdateWaitingReservationNote(ctx context.Context, req *guestProto.UpdateWaitingReservationNoteRequest) (*guestProto.UpdateWaitingReservationNoteResponse, error) {
	res, err := s.reservationWaitListService.UpdateWaitingReservationNote(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
