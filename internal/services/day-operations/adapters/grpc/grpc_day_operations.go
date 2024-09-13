package grpc

import "github.com/goplaceapp/goplace-guest/internal/services/day-operations/application"

type DayOperationsServiceServer struct {
	dayOperationsService *application.DayOperationsService
}

func NewDayOperationsServiceServer() *DayOperationsServiceServer {
	return &DayOperationsServiceServer{
		dayOperationsService: application.NewDayOperationsService(),
	}
}
