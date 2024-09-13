package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s*ReservationFeedbackServiceServer) GetAllReservationFeedbackSections(ctx context.Context, req*emptypb.Empty)(*guestProto.GetAllReservationsFeedbackSectionsResponse,error){
	res, err := s.reservationFeedbackService.GetAllReservationFeedbackSections(ctx,req)
	if err != nil{
		return nil,err
	}
	return res,nil
}