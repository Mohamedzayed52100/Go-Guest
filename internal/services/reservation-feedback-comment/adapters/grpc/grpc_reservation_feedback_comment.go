package grpc

import "github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback-comment/application"

type ReservationFeedbackCommentServiceServer struct {
	reservationFeedbackCommentService *application.ReservationFeedbackCommentService
}

func NewReservationFeedbackCommentServiceServer() *ReservationFeedbackCommentServiceServer {
	return &ReservationFeedbackCommentServiceServer{
		reservationFeedbackCommentService: application.NewReservationFeedbackCommentService(),
	}
}
