package application

import "github.com/goplaceapp/goplace-guest/internal/services/payment/infrastructure/repository"

type PaymentService struct {
	Repository *repository.PaymentRepository
}

func NewPaymentService() *PaymentService {
	return &PaymentService{
		Repository: repository.NewPaymentRepository(),
	}
}
