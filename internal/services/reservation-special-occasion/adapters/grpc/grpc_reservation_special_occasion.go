package grpc

import "github.com/goplaceapp/goplace-guest/internal/services/reservation-special-occasion/application"

type ReservationSpecialOccasionServiceServer struct {
	reservationSpecialOccasionService *application.ReservationSpecialOccasionService
}

func NewReservationSpecialOccasionServiceServer() *ReservationSpecialOccasionServiceServer {
	return &ReservationSpecialOccasionServiceServer{
		reservationSpecialOccasionService: application.NewReservationSpecialOccasionService(),
	}
}
