package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationWaitListServiceServer) GetAllWaitingReservations(ctx context.Context, req *guestProto.GetWaitingReservationRequest) (*guestProto.GetWaitingReservationsResponse, error) {
	res, err := s.reservationWaitListService.GetAllWaitingReservations(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *ReservationWaitListServiceServer) GetRealtimeWaitingReservations(req *guestProto.GetWaitingReservationRequest, stream guestProto.ReservationWaitlist_GetRealtimeWaitingReservationsServer) error {
	err := s.reservationWaitListService.GetRealtimeWaitingReservations(req, stream)
	if err != nil {
		return err
	}
	return nil
}
