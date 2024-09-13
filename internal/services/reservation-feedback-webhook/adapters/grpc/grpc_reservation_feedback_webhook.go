package grpc

import "github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback-webhook/application"

type ReservationFeedbackWebhookServiceServer struct {
	reservationFeedbackWebhookService *application.ReservationFeedbackWebhookService
}

func NewReservationFeedbackWebhookServiceServer() *ReservationFeedbackWebhookServiceServer {
	return &ReservationFeedbackWebhookServiceServer{
		reservationFeedbackWebhookService: application.NewReservationFeedbackWebhookService(),
	}
}
