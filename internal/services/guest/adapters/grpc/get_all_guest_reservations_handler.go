package grpc

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestServiceServer) GetAllGuestReservations(ctx context.Context, req *guestProto.GetAllGuestReservationsRequest) (*guestProto.GetAllGuestReservationsResponse, error) {
	res, err := s.guestService.GetAllGuestReservations(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
