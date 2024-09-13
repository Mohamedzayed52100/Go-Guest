package grpc

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestServiceServer) GetGuestSpending(ctx context.Context, req *guestProto.GetGuestSpendingRequest) (*guestProto.GetGuestSpendingResponse, error) {
	res, err := s.guestService.GetGuestSpending(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
