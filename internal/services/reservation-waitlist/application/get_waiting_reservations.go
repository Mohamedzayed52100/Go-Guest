package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationWaitListService) GetAllWaitingReservations(ctx context.Context, req *guestProto.GetWaitingReservationRequest) (*guestProto.GetWaitingReservationsResponse, error) {
	res, err := s.Repository.GetAllWaitingReservations(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *ReservationWaitListService) GetRealtimeWaitingReservations(req *guestProto.GetWaitingReservationRequest, stream guestProto.ReservationWaitlist_GetRealtimeWaitingReservationsServer) error {
	err := s.Repository.GetRealtimeWaitingReservations(req,stream)
	if err != nil {
		return err
	}
	return nil
}
