package application

import "github.com/goplaceapp/goplace-guest/internal/services/reservation-log/infrastructure/repository"

type ReservationLogService struct {
	Repository *repository.ReservationLogRepository
}

func NewReservationLogService() *ReservationLogService {
	return &ReservationLogService{
		Repository: repository.NewReservationLogRepository(),
	}
}
