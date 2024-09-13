package grpc

import "github.com/goplaceapp/goplace-guest/internal/services/payment/application"

type PaymentServiceServer struct {
	paymentService *application.PaymentService
}

func NewPaymentServiceServer() *PaymentServiceServer {
	return &PaymentServiceServer{
		paymentService: application.NewPaymentService(),
	}
}
