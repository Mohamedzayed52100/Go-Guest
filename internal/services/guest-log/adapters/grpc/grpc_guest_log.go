package grpc

import "github.com/goplaceapp/goplace-guest/internal/services/guest-log/application"

type GuestLogServiceServer struct {
	guestLogService *application.GuestLogService
}

func NewGuestLogServiceServer() *GuestLogServiceServer {
	return &GuestLogServiceServer{
		guestLogService: application.NewGuestLogService(),
	}
}
