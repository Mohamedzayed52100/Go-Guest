package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationWaitListService) CreateDirectInReservation(ctx context.Context, req *guestProto.CreateWaitingReservationRequest) (*guestProto.CreateReservationResponse, error) {
	res, err := s.Repository.CreateDirectInReservation(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
