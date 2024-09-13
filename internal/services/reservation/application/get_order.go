package application

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationService) GetReservationOrderByReservationID(ctx context.Context, req *guestProto.GetReservationOrderByReservationIDRequest) (*guestProto.GetReservationOrderByReservationIDResponse, error) {
	res, err := s.Repository.GetReservationOrderByReservationID(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
