package application

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestService) GetGuestSpending(ctx context.Context, req *guestProto.GetGuestSpendingRequest) (*guestProto.GetGuestSpendingResponse, error) {
	res, err := s.Repository.GetGuestSpending(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
