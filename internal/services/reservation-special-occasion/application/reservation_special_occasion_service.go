package application

import "github.com/goplaceapp/goplace-guest/internal/services/reservation-special-occasion/infrastructure/repository"

type ReservationSpecialOccasionService struct {
	Repository *repository.ReservationSpecialOccasionRepository
}

func NewReservationSpecialOccasionService() *ReservationSpecialOccasionService {
	return &ReservationSpecialOccasionService{
		Repository: repository.NewReservationSpecialOccasionRepository(),
	}
}
