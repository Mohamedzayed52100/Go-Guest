package application

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationService) DeleteReservation(ctx context.Context, req *guestProto.DeleteReservationRequest) (*guestProto.DeleteReservationResponse, error) {
	res, err := s.Repository.DeleteReservation(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
