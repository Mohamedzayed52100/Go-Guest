package application

import "github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback-comment/infrastructure/repository"

type ReservationFeedbackCommentService struct {
	Repository *repository.ReservationFeedbackCommentRepository
}

func NewReservationFeedbackCommentService() *ReservationFeedbackCommentService {
	return &ReservationFeedbackCommentService{
		Repository: repository.NewReservationFeedbackCommentRepository(),
	}
}
