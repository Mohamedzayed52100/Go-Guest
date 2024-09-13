package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationFeedbackCommentService) GetAllReservationFeedbackComments(ctx context.Context, req *guestProto.GetAllReservationFeedbackCommentsRequest) (*guestProto.GetAllReservationFeedbackCommentsResponse, error) {
	res, err := s.Repository.GetAllReservationFeedbackComments(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
