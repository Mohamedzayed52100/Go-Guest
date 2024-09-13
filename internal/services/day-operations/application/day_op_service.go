package application

import "github.com/goplaceapp/goplace-guest/internal/services/day-operations/infrastructure/repository"

type DayOperationsService struct {
	Repository *repository.DayOperationsRepository
}

func NewDayOperationsService() *DayOperationsService {
	return &DayOperationsService{
		Repository: repository.NewDayOperationsRepository(),
	}
}
