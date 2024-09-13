package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestServiceServer) DeleteReservationVisitors(ctx context.Context, req *guestProto.DeleteReservationVisitorsRequest) (*guestProto.DeleteReservationVisitorsResponse, error) {
	res, err := s.guestService.DeleteReservationVisitors(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
