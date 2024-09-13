package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationLogServiceServer) GetAllReservationLogs(ctx context.Context, req *guestProto.GetAllReservationLogsRequest) (*guestProto.GetAllReservationLogsResponse, error) {
	res, err := s.reservationLogService.GetAllReservationLogs(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, err
}
