package application

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestService) GetGuestUpcomingReservations(ctx context.Context, req *guestProto.GetGuestUpcomingReservationsRequest) (*guestProto.GetGuestUpcomingReservationsResponse, error) {
	res, err := s.Repository.GetGuestUpcomingReservations(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
