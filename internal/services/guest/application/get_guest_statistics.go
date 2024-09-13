package application

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestService) GetGuestStatistics(ctx context.Context, req *guestProto.GetGuestStatisticsRequest) (*guestProto.GetGuestStatisticsResponse, error) {
	res, err := s.Repository.GetGuestStatistics(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
