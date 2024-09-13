package grpc

import "github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback/application"

type ReservationFeedbackServiceServer struct {
	reservationFeedbackService *application.ReservationFeedbackService
}

func NewReservationFeedbackServiceServer() *ReservationFeedbackServiceServer {
	return &ReservationFeedbackServiceServer{
		reservationFeedbackService: application.NewReservationFeedbackService(),
	}
}
