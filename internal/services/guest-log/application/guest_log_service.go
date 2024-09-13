package application

import "github.com/goplaceapp/goplace-guest/internal/services/guest-log/infrastructure/repository"

type GuestLogService struct {
	Repository *repository.GuestLogRepository
}

func NewGuestLogService() *GuestLogService {
	return &GuestLogService{
		Repository: repository.NewGuestLogRepository(),
	}
}
