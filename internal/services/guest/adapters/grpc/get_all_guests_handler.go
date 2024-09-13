package grpc

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestServiceServer) GetAllGuests(ctx context.Context, req *guestProto.GetAllGuestsRequest) (*guestProto.GetAllGuestsResponse, error) {
	res, err := s.guestService.GetAllGuests(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
