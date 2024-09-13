package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationServiceServer) AddReservationNote(ctx context.Context, req *guestProto.AddReservationNoteRequest) (*guestProto.AddReservationNoteResponse, error) {
	res, err := s.reservationService.AddReservationNote(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
