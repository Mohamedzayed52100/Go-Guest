package main

import (
	"os"

	"github.com/goplaceapp/goplace-common/pkg/grpchelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-common/pkg/utils"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/server"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	defer utils.RecoverFromPanic()

	godotenv.Load()
	meta.TokenSymmetricKey = os.Getenv("TOKEN_SYMMETRIC_KEY")
	service := server.New()
	grpchelper.Start(func(server *grpc.Server) {
		guestProto.RegisterGuestServer(server, service.GuestServiceServer)
		guestProto.RegisterGuestLogServer(server, service.GuestLogServiceServer)
		guestProto.RegisterReservationServer(server, service.ReservationServiceServer)
		guestProto.RegisterReservationLogServer(server, service.ReservationLogServiceServer)
		guestProto.RegisterReservationSpecialOccasionServer(server, service.ReservationSpecialOccasionServiceServer)
		guestProto.RegisterReservationFeedbackServer(server, service.ReservationFeedbackServiceServer)
		guestProto.RegisterReservationFeedbackCommentServer(server, service.ReservationFeedbackCommentServiceServer)
		guestProto.RegisterReservationWaitlistServer(server, service.ReservationWaitListServiceServer)
		guestProto.RegisterReservationFeedbackWebhookServer(server, service.ReservationFeedbackWebhookServiceServer)
		guestProto.RegisterDayOperationsServer(server, service.DayOperationsServiceServer)
		guestProto.RegisterReservationWidgetServer(server, service.ReservationWidgetServiceServer)
		guestProto.RegisterPaymentServer(server, service.PaymentServiceServer)
	}, service)
}
