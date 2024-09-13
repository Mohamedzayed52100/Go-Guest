package grpc

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestServiceServer) GetGuestStatistics(ctx context.Context, req *guestProto.GetGuestStatisticsRequest) (*guestProto.GetGuestStatisticsResponse, error) {
	res, err := s.guestService.GetGuestStatistics(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
