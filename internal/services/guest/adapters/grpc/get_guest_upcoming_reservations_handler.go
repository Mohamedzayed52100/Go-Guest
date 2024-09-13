package grpc

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestServiceServer) GetGuestUpcomingReservations(ctx context.Context, req *guestProto.GetGuestUpcomingReservationsRequest) (*guestProto.GetGuestUpcomingReservationsResponse, error) {
	res, err := s.guestService.GetGuestUpcomingReservations(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
