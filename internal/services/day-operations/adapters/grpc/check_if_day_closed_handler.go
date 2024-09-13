package grpc

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *DayOperationsServiceServer) CheckIfDayClosed(ctx context.Context, req *guestProto.CheckIfDayClosedRequest) (*guestProto.CheckIfDayClosedResponse, error) {
	res, err := s.dayOperationsService.CheckIfDayClosed(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
