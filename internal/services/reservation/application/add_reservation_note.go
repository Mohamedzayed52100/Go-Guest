package application

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationService) AddReservationNote(ctx context.Context, req *guestProto.AddReservationNoteRequest) (*guestProto.AddReservationNoteResponse, error) {
	res, err := s.Repository.AddReservationNote(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
