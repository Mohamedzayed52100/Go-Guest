package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestService) DeleteReservationVisitors(ctx context.Context, req *guestProto.DeleteReservationVisitorsRequest) (*guestProto.DeleteReservationVisitorsResponse, error) {
	res, err := s.Repository.DeleteReservationVisitors(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
