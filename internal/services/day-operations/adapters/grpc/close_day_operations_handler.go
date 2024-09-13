package grpc

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *DayOperationsServiceServer) CloseDayOperations(ctx context.Context, req *guestProto.CloseDayOperationsRequest) (*guestProto.CloseDayOperationsResponse, error) {
	res, err := s.dayOperationsService.CloseDayOperations(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
