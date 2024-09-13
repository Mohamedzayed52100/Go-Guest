package application

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationService) GetAvailableTimes(ctx context.Context, req *guestProto.GetAvailableTimesRequest) (*guestProto.GetAvailableTimesResponse, error) {
	res, err := s.Repository.GetAvailableTimes(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
