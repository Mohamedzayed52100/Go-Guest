package grpc

import (
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"golang.org/x/net/context"
)

func (s *GuestServiceServer) GetGuestReservationStatistics(ctx context.Context, req *guestProto.GetGuestReservationStatisticsRequest) (*guestProto.GetGuestReservationStatisticsResponse, error) {
	res, err := s.guestService.GetGuestReservationStatistics(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
