package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s*ReservationFeedbackServiceServer) UpdateReservationFeedback(ctx context.Context, req*guestProto.UpdateReservationFeedbackRequest)(*guestProto.UpdateReservationFeedbackResponse,error){
	res, err := s.reservationFeedbackService.UpdateReservationFeedback(ctx,req)
	if err != nil{
		return nil,err
	}
	return res,nil
}