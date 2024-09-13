package grpc

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestServiceServer) GetGuestByID(ctx context.Context, req *guestProto.GetGuestByIDRequest) (*guestProto.GetGuestByIDResponse, error) {
	res, err := s.guestService.GetGuestByID(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
