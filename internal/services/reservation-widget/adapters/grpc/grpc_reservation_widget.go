package grpc

import "github.com/goplaceapp/goplace-guest/internal/services/reservation-widget/application"

type ReservationWidgetServiceServer struct {
	reservationWidgetService *application.ReservationWidgetService
}

func NewReservationServiceServer() *ReservationWidgetServiceServer {
	return &ReservationWidgetServiceServer{
		reservationWidgetService: application.NewReservationWidgetService(),
	}
}
