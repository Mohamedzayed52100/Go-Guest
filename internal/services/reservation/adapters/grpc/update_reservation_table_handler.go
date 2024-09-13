package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationServiceServer) UpdateReservationTable(ctx context.Context, req *guestProto.UpdateReservationTableRequest) (*guestProto.UpdateReservationTableResponse, error) {
	res, err := s.reservationService.UpdateReservationTable(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
