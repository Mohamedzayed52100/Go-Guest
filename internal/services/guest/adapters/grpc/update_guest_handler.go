package grpc

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestServiceServer) UpdateGuest(ctx context.Context, req *guestProto.UpdateGuestRequest) (*guestProto.UpdateGuestResponse, error) {
	res, err := s.guestService.UpdateGuest(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
