package grpc

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestServiceServer) CreateGuest(ctx context.Context, req *guestProto.CreateGuestRequest) (*guestProto.CreateGuestResponse, error) {
	res, err := s.guestService.CreateGuest(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
