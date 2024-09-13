package application

import "github.com/goplaceapp/goplace-guest/internal/services/guest/infrastructure/repository"

type GuestService struct {
	Repository *repository.GuestRepository
}

func NewGuestService() *GuestService {
	return &GuestService{
		Repository: repository.NewGuestRepository(),
	}
}
