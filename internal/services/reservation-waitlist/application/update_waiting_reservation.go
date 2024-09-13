package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationWaitListService) UpdateWaitingReservationDetails(ctx context.Context, req *guestProto.UpdateWaitingReservationDetailsRequest) (*guestProto.UpdateWaitingReservationDetailsResponse, error) {
	res, err := s.Repository.UpdateWaitingReservation(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
