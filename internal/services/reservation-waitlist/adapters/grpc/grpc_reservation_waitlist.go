package grpc

import "github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/application"

type ReservationWaitListServiceServer struct {
	reservationWaitListService *application.ReservationWaitListService
}

func NewReservationWaitListServiceServer() *ReservationWaitListServiceServer {
	return &ReservationWaitListServiceServer{
		reservationWaitListService: application.NewReservationWaitListService(),
	}
}
