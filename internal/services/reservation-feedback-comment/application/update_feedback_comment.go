package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationFeedbackCommentService) UpdateReservationFeedbackComment(ctx context.Context, req *guestProto.UpdateReservationFeedbackCommentRequest) (*guestProto.UpdateReservationFeedbackCommentResponse, error) {
	res, err := s.Repository.UpdateReservationFeedbackComment(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
