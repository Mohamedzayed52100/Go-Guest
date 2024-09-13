package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestLogService) GetAllGuestLogs(ctx context.Context, req *guestProto.GetAllGuestLogsRequest) (*guestProto.GetAllGuestLogsResponse, error) {
	res, err := s.Repository.GetAllGuestLogs(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, err
}
