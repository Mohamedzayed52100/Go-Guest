package application

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationService) CreateReservation(ctx context.Context, req *guestProto.CreateReservationRequest) (*guestProto.CreateReservationResponse, error) {
	res, err := s.Repository.CreateReservation(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
