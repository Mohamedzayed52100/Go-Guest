package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s*ReservationServiceServer) GetReservationOrderByReservationID(ctx context.Context, req*guestProto.GetReservationOrderByReservationIDRequest)(*guestProto.GetReservationOrderByReservationIDResponse,error){
	res, err := s.reservationService.GetReservationOrderByReservationID(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}