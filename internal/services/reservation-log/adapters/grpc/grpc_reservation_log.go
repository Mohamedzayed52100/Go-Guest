package grpc

import "github.com/goplaceapp/goplace-guest/internal/services/reservation-log/application"

type ReservationLogServiceServer struct {
	reservationLogService *application.ReservationLogService
}

func NewReservationLogServiceServer() *ReservationLogServiceServer {
	return &ReservationLogServiceServer{
		reservationLogService: application.NewReservationLogService(),
	}
}
