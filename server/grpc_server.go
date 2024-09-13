package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/goplaceapp/goplace-common/pkg/grpchelper"
	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/goplaceapp/goplace-guest/database"
	dayOperationsAdapters "github.com/goplaceapp/goplace-guest/internal/services/day-operations/adapters/grpc"
	guestLogAdapters "github.com/goplaceapp/goplace-guest/internal/services/guest-log/adapters/grpc"
	guestAdapters "github.com/goplaceapp/goplace-guest/internal/services/guest/adapters/grpc"
	paymentAdapters "github.com/goplaceapp/goplace-guest/internal/services/payment/adapters/grpc"
	reservationFeedbackCommentAdapters "github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback-comment/adapters/grpc"
	reservationFeedbackWebhookAdapters "github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback-webhook/adapters/grpc"
	reservationFeedbackAdapters "github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback/adapters/grpc"
	reservationLogAdapters "github.com/goplaceapp/goplace-guest/internal/services/reservation-log/adapters/grpc"
	reservationSpecialOccasionAdapters "github.com/goplaceapp/goplace-guest/internal/services/reservation-special-occasion/adapters/grpc"
	reservationWaitListAdapters "github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/adapters/grpc"
	reservationWidgetAdapters "github.com/goplaceapp/goplace-guest/internal/services/reservation-widget/adapters/grpc"
	reservationAdapters "github.com/goplaceapp/goplace-guest/internal/services/reservation/adapters/grpc"
	"go.uber.org/zap"
)

type Service struct {
	Log                                     *zap.SugaredLogger
	BaseCfg                                 *grpchelper.BaseConfig
	PostgresService                         *database.PostgresService
	GuestServiceServer                      *guestAdapters.GuestServiceServer
	GuestLogServiceServer                   *guestLogAdapters.GuestLogServiceServer
	ReservationServiceServer                *reservationAdapters.ReservationServiceServer
	DayOperationsServiceServer              *dayOperationsAdapters.DayOperationsServiceServer
	ReservationLogServiceServer             *reservationLogAdapters.ReservationLogServiceServer
	ReservationFeedbackServiceServer        *reservationFeedbackAdapters.ReservationFeedbackServiceServer
	ReservationWaitListServiceServer        *reservationWaitListAdapters.ReservationWaitListServiceServer
	ReservationFeedbackCommentServiceServer *reservationFeedbackCommentAdapters.ReservationFeedbackCommentServiceServer
	ReservationFeedbackWebhookServiceServer *reservationFeedbackWebhookAdapters.ReservationFeedbackWebhookServiceServer
	ReservationSpecialOccasionServiceServer *reservationSpecialOccasionAdapters.ReservationSpecialOccasionServiceServer
	ReservationWidgetServiceServer          *reservationWidgetAdapters.ReservationWidgetServiceServer
	PaymentServiceServer                    *paymentAdapters.PaymentServiceServer
}

func New() *Service {
	// initialize logger
	log, err := logger.New(os.Getenv("LOG_LEVEL"))
	if err != nil {
		panic(fmt.Errorf("failed to initialize the logger, %w", err))
	}

	// Create an instance of the InquiryCommentServiceServer
	postgresService := database.NewService()

	// Create an instance of the GuestServiceServer
	guestServiceServer := guestAdapters.NewGuestServiceServer()

	// Create an instance of the GuestLogServiceServer
	guestLogServiceServer := guestLogAdapters.NewGuestLogServiceServer()

	// Create an instance of the ReservationServiceServer
	reservationServiceServer := reservationAdapters.NewReservationServiceServer()

	// Create an instance of the ReservationWaitListServiceServer
	reservationWaitListServiceServer := reservationWaitListAdapters.NewReservationWaitListServiceServer()

	// Create an instance of the ReservationLogServiceServer
	reservationLogServiceServer := reservationLogAdapters.NewReservationLogServiceServer()

	// Create an instance of the ReservationFeedbackServiceServer
	reservationFeedbackServiceServer := reservationFeedbackAdapters.NewReservationFeedbackServiceServer()

	// Create an instance of the ReservationWidgetServiceServer
	reservationWidgetServiceServer := reservationWidgetAdapters.NewReservationServiceServer()

	// Create an instance of the ReservationFeedbackCommentServiceServer
	reservationFeedbackCommentServiceServer := reservationFeedbackCommentAdapters.NewReservationFeedbackCommentServiceServer()

	// Create an instance of the ReservationFeedbackWebhookServiceServer
	reservationFeedbackWebhookServiceServer := reservationFeedbackWebhookAdapters.NewReservationFeedbackWebhookServiceServer()

	// Create an instance of the ReservationSpecialOccasionServiceServer
	reservationSpecialOccasionServiceServer := reservationSpecialOccasionAdapters.NewReservationSpecialOccasionServiceServer()

	// Create an instance of the DayOperationsServiceServer
	dayOperationsServiceServer := dayOperationsAdapters.NewDayOperationsServiceServer()

	paymentServer := paymentAdapters.NewPaymentServiceServer()

	return &Service{
		Log:                                     log,
		BaseCfg:                                 &grpchelper.BaseConfig{},
		PostgresService:                         postgresService,
		GuestServiceServer:                      guestServiceServer,
		GuestLogServiceServer:                   guestLogServiceServer,
		ReservationServiceServer:                reservationServiceServer,
		ReservationWaitListServiceServer:        reservationWaitListServiceServer,
		ReservationLogServiceServer:             reservationLogServiceServer,
		ReservationFeedbackServiceServer:        reservationFeedbackServiceServer,
		ReservationFeedbackCommentServiceServer: reservationFeedbackCommentServiceServer,
		ReservationFeedbackWebhookServiceServer: reservationFeedbackWebhookServiceServer,
		ReservationSpecialOccasionServiceServer: reservationSpecialOccasionServiceServer,
		DayOperationsServiceServer:              dayOperationsServiceServer,
		ReservationWidgetServiceServer:          reservationWidgetServiceServer,
		PaymentServiceServer:                    paymentServer,
	}
}

func (s *Service) SetBaseConfig(cfg *grpchelper.BaseConfig) {
	s.BaseCfg = cfg
}

func (s *Service) GetLog() *zap.SugaredLogger {
	return s.Log
}

func (s *Service) GetLivenessHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		// service specific health state definition goes here
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Service) GetReadinessHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		// service specific ready state definition goes here
		w.WriteHeader(http.StatusOK)
	}
}
