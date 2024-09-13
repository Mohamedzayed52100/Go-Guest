package application

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *DayOperationsService) CheckIfDayClosed(ctx context.Context, req *guestProto.CheckIfDayClosedRequest) (*guestProto.CheckIfDayClosedResponse, error) {
	res, err := s.Repository.CheckIfDayClosed(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
