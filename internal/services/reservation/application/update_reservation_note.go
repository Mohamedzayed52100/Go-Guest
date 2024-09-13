package application

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationService) UpdateReservationNote(ctx context.Context, req *guestProto.UpdateReservationNoteRequest) (*guestProto.UpdateReservationNoteResponse, error) {
	res, err := s.Repository.UpdateReservationNote(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
