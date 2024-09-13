package grpc

import "github.com/goplaceapp/goplace-guest/internal/services/reservation/application"

type ReservationServiceServer struct {
	reservationService *application.ReservationService
}

func NewReservationServiceServer() *ReservationServiceServer {
	return &ReservationServiceServer{
		reservationService: application.NewReservationService(),
	}
}
