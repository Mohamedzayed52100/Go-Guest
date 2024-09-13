package grpc

import "github.com/goplaceapp/goplace-guest/internal/services/guest/application"

type GuestServiceServer struct {
	guestService *application.GuestService
}

func NewGuestServiceServer() *GuestServiceServer {
	return &GuestServiceServer{
		guestService: application.NewGuestService(),
	}
}
