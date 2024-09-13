package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationWidgetService) CreateWidgetReservation(ctx context.Context, req *guestProto.CreateWidgetReservationRequest) (*guestProto.CreateWidgetReservationResponse, error) {
	res, err := s.Repository.CreateWidgetReservation(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
