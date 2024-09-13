package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s*ReservationFeedbackService) GetAllReservationFeedbackSections(ctx context.Context, req*emptypb.Empty)(*guestProto.GetAllReservationsFeedbackSectionsResponse,error){
	res, err := s.Repository.GetAllReservationFeedbackSections(ctx,req)
	if err != nil{
		return nil,err
	}
	return res,nil
}