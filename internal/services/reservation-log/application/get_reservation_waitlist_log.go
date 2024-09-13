package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationLogService) GetAllReservationWaitlistLogs(ctx context.Context, req *guestProto.GetAllReservationWaitlistLogsRequest) (*guestProto.GetAllReservationWaitlistLogsResponse, error) {
	res, err := s.Repository.GetAllReservationWaitlistLogs(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, err
}
