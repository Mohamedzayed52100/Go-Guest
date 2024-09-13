package application

import (
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"golang.org/x/net/context"
)

func (s *GuestService) GetGuestReservationStatistics(ctx context.Context, req *guestProto.GetGuestReservationStatisticsRequest) (*guestProto.GetGuestReservationStatisticsResponse, error) {
	res, err := s.Repository.GetGuestReservationStatistics(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
