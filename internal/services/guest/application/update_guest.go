package application

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestService) UpdateGuest(ctx context.Context, req *guestProto.UpdateGuestRequest) (*guestProto.UpdateGuestResponse, error) {
	res, err := s.Repository.UpdateGuest(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
