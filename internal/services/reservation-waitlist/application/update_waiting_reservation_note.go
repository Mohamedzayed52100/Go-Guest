package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationWaitListService) UpdateWaitingReservationNote(ctx context.Context, req *guestProto.UpdateWaitingReservationNoteRequest) (*guestProto.UpdateWaitingReservationNoteResponse, error) {
	res, err := s.Repository.UpdateWaitingReservationNote(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
