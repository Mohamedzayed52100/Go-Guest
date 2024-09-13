package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationServiceServer) GetAllReservations(ctx context.Context, req *guestProto.GetAllReservationsRequest) (*guestProto.GetAllReservationsResponse, error) {
	res, err := s.reservationService.GetAllReservations(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *ReservationServiceServer) GetReservationByID(ctx context.Context, req *guestProto.GetReservationByIDRequest) (*guestProto.GetReservationByIDResponse, error) {
	res, err := s.reservationService.GetReservationByID(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *ReservationServiceServer) GetRealtimeReservations(req *guestProto.GetRealtimeReservationsRequest, stream guestProto.Reservation_GetRealtimeReservationsServer) error {
	err := s.reservationService.GetRealtimeReservations(req, stream)
	if err != nil {
		return err
	}
	return nil
}

