package grpc

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/goplaceapp/goplace-guest/config"
	"github.com/goplaceapp/goplace-guest/database"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation/application"
	"gorm.io/gorm"
)

type ReservationServiceServer struct {
	ReservationService *application.ReservationService
}

func NewReservationServiceServer(db *gorm.DB, tenantDBConnections map[string]*gorm.DB) *ReservationServiceServer {
	if database.SharedPostgresService == nil {
		cfg := &config.Config{}
		if err := env.Parse(cfg); err != nil {
			panic(fmt.Errorf("failed to parse environment variables, %w", err))
		}

		database.SharedPostgresService = &database.PostgresService{
			Db:                  db,
			TenantDbConnections: tenantDBConnections,
			SvcCfg:              cfg,
		}
	}

	return &ReservationServiceServer{
		ReservationService: application.NewReservationService(),
	}
}
