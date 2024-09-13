package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s*ReservationFeedbackService) CreateReservationFeedback(ctx context.Context, req*guestProto.CreateReservationFeedbackRequest)(*guestProto.CreateReservationFeedbackResponse,error){
	res, err := s.Repository.CreateReservationFeedback(ctx,req)
	if err != nil{
		return nil,err
	}
	return res,nil
}