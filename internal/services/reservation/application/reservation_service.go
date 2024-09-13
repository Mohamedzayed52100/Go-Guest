package application

import "github.com/goplaceapp/goplace-guest/internal/services/reservation/infrastructure/repository"

type ReservationService struct {
	Repository *repository.ReservationRepository
}

func NewReservationService() *ReservationService {
	return &ReservationService{
		Repository: repository.NewReservationRepository(),
	}
}
