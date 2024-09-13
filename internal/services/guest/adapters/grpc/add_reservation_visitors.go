package grpc

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestServiceServer) AddReservationVisitors(ctx context.Context, req *guestProto.AddReservationVisitorsRequest) (*guestProto.AddReservationVisitorsResponse, error) {
	res, err := s.guestService.AddReservationVisitors(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
