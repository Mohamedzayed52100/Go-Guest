package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationService) UpdateReservationTable(ctx context.Context, req *guestProto.UpdateReservationTableRequest) (*guestProto.UpdateReservationTableResponse, error) {
	res, err := s.Repository.UpdateReservationTable(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
