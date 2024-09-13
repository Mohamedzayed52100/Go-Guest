package grpc

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationWidgetServiceServer) GetWidgetAvailableTimes(ctx context.Context, req *guestProto.GetWidgetAvailableTimesRequest) (*guestProto.GetWidgetAvailableTimesResponse, error) {
	res, err := s.reservationWidgetService.GetWidgetAvailableTimes(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
