package application

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *DayOperationsService) CloseDayOperations(ctx context.Context, req *guestProto.CloseDayOperationsRequest) (*guestProto.CloseDayOperationsResponse, error) {
	res, err := s.Repository.CloseDayOperations(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
