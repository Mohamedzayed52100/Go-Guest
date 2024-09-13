package application

import "github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback/infrastructure/repository"

type ReservationFeedbackService struct {
	Repository *repository.ReservationFeedbackRepository
}

func NewReservationFeedbackService() *ReservationFeedbackService {
	return &ReservationFeedbackService{
		Repository: repository.NewReservationFeedbackRepository(),
	}
}
