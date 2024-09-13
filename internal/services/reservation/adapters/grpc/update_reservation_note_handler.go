package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationServiceServer) UpdateReservationNote(ctx context.Context, req *guestProto.UpdateReservationNoteRequest) (*guestProto.UpdateReservationNoteResponse, error) {
	res, err := s.reservationService.UpdateReservationNote(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
