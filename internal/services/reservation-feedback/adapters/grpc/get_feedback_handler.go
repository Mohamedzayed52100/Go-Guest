package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationFeedbackServiceServer) GetAllReservationsFeedbacks(ctx context.Context, req *guestProto.GetAllReservationsFeedbacksRequest) (*guestProto.GetAllReservationsFeedbacksResponse, error) {
	res, err := s.reservationFeedbackService.GetAllReservationsFeedbacks(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *ReservationFeedbackServiceServer) GetReservationFeedbackByID(ctx context.Context, req *guestProto.GetReservationFeedbackByIDRequest) (*guestProto.GetReservationFeedbackByIDResponse, error) {
	res, err := s.reservationFeedbackService.GetReservationFeedbackByID(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}