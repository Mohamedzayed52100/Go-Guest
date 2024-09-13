package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestService) GetGuestByID(ctx context.Context, req *guestProto.GetGuestByIDRequest) (*guestProto.GetGuestByIDResponse, error) {
	res, err := s.Repository.CommonRepository.GetGuestByID(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
