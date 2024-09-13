package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationFeedbackCommentServiceServer) AddReservationFeedbackSolution(ctx context.Context, req *guestProto.AddReservationFeedbackSolutionRequest) (*guestProto.AddReservationFeedbackSolutionResponse, error) {
	res, err := s.reservationFeedbackCommentService.AddReservationFeedbackSolution(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
