package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationFeedbackCommentServiceServer) GetAllReservationFeedbackComments(ctx context.Context, req *guestProto.GetAllReservationFeedbackCommentsRequest) (*guestProto.GetAllReservationFeedbackCommentsResponse, error) {
	res, err := s.reservationFeedbackCommentService.GetAllReservationFeedbackComments(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
