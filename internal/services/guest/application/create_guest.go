package application

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestService) CreateGuest(ctx context.Context, req *guestProto.CreateGuestRequest) (*guestProto.CreateGuestResponse, error) {
	res, err := s.Repository.CreateGuest(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
