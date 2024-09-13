package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s*ReservationFeedbackServiceServer) CreateReservationFeedback(ctx context.Context, req*guestProto.CreateReservationFeedbackRequest)(*guestProto.CreateReservationFeedbackResponse,error){
	res, err := s.reservationFeedbackService.CreateReservationFeedback(ctx,req)
	if err != nil{
		return nil,err
	}
	return res,nil
}