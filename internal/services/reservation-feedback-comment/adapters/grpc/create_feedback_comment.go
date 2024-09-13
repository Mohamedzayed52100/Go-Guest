package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationFeedbackCommentServiceServer) CreateReservationFeedbackComment(ctx context.Context, req *guestProto.CreateReservationFeedbackCommentRequest) (*guestProto.CreateReservationFeedbackCommentResponse, error) {
	res, err := s.reservationFeedbackCommentService.CreateReservationFeedbackComment(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
