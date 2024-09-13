package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationFeedbackService) GetAllReservationsFeedbacks(ctx context.Context, req *guestProto.GetAllReservationsFeedbacksRequest) (*guestProto.GetAllReservationsFeedbacksResponse, error) {
	res, err := s.Repository.GetAllReservationsFeedbacks(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *ReservationFeedbackService) GetReservationFeedbackByID(ctx context.Context, req *guestProto.GetReservationFeedbackByIDRequest) (*guestProto.GetReservationFeedbackByIDResponse, error) {
	res, err := s.Repository.GetReservationFeedbackByID(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
