package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationWaitListService) CreateWaitingReservationNote(ctx context.Context, req *guestProto.CreateWaitingReservationNoteRequest) (*guestProto.CreateWaitingReservationNoteResponse, error) {
	res, err := s.Repository.CreateWaitingReservationNote(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
