package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationWaitListService) CreateWaitingReservation(ctx context.Context, req *guestProto.CreateWaitingReservationRequest) (*guestProto.CreateWaitingReservationResponse, error) {
	res, err := s.Repository.CreateWaitingReservation(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
