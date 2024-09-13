package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationService) GetAllReservations(ctx context.Context, req *guestProto.GetAllReservationsRequest) (*guestProto.GetAllReservationsResponse, error) {
	res, err := s.Repository.GetAllReservations(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *ReservationService) GetReservationByID(ctx context.Context, req *guestProto.GetReservationByIDRequest) (*guestProto.GetReservationByIDResponse, error) {
	res, err := s.Repository.GetReservationByID(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *ReservationService) GetRealtimeReservations(req *guestProto.GetRealtimeReservationsRequest, stream guestProto.Reservation_GetRealtimeReservationsServer)  error {
	err := s.Repository.GetRealTimeReservations(req, stream)
	if err != nil {
		return err
	}
	return nil
}
