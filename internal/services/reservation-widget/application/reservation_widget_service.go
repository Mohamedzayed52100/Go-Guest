package application

import "github.com/goplaceapp/goplace-guest/internal/services/reservation-widget/infrastructure/repository"

type ReservationWidgetService struct {
	Repository *repository.ReservationWidgetRepository
}

func NewReservationWidgetService() *ReservationWidgetService {
	return &ReservationWidgetService{
		Repository: repository.NewReservationWidgetRepository(),
	}
}
