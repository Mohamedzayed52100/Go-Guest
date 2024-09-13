package application

import (
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"golang.org/x/net/context"
)

func (s *ReservationWidgetService) GetWidgetAvailableTimes(ctx context.Context, req *guestProto.GetWidgetAvailableTimesRequest) (*guestProto.GetWidgetAvailableTimesResponse, error) {
	res, err := s.Repository.GetWidgetAvailableTimes(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
