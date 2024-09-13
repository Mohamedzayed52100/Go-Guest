package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s*ReservationFeedbackService) UpdateReservationFeedback(ctx context.Context, req*guestProto.UpdateReservationFeedbackRequest)(*guestProto.UpdateReservationFeedbackResponse,error){
	res, err := s.Repository.UpdateReservationFeedback(ctx,req)
	if err != nil{
		return nil,err
	}
	return res,nil
}