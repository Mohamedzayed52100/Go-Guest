package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestLogServiceServer) GetAllGuestLogs(ctx context.Context, req *guestProto.GetAllGuestLogsRequest) (*guestProto.GetAllGuestLogsResponse, error) {
	res, err := s.guestLogService.GetAllGuestLogs(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, err
}
