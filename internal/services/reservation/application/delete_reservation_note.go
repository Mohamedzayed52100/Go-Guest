package application

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationService) DeleteReservationNote(ctx context.Context, req *guestProto.DeleteReservationNoteRequest) (*guestProto.DeleteReservationNoteResponse, error) {
	res, err := s.Repository.DeleteReservationNote(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
