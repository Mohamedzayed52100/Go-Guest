package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationLogService) GetAllReservationLogs(ctx context.Context, req *guestProto.GetAllReservationLogsRequest) (*guestProto.GetAllReservationLogsResponse, error) {
	res, err := s.Repository.GetAllReservationLogs(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, err
}
