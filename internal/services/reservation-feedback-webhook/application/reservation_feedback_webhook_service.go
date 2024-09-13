package application

import "github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback-webhook/infrastructure/repository"

type ReservationFeedbackWebhookService struct {
	Repository *repository.ReservationFeedbackWebhookRepository
}

func NewReservationFeedbackWebhookService() *ReservationFeedbackWebhookService {
	return &ReservationFeedbackWebhookService{
		Repository: repository.NewReservationFeedbackWebhookRepository(),
	}
}
