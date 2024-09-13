package application

import "github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/infrastructure/repository"

type ReservationWaitListService struct {
	Repository *repository.ReservationWaitListRepository
}

func NewReservationWaitListService() *ReservationWaitListService {
	return &ReservationWaitListService{
		Repository: repository.NewReservationWaitListRepository(),
	}
}
