package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationFeedbackCommentService) AddReservationFeedbackSolution(ctx context.Context, req *guestProto.AddReservationFeedbackSolutionRequest) (*guestProto.AddReservationFeedbackSolutionResponse, error) {
	res, err := s.Repository.AddReservationFeedbackSolution(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
