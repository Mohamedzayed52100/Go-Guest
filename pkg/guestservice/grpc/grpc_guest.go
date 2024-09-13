package grpc

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/goplaceapp/goplace-guest/config"
	"github.com/goplaceapp/goplace-guest/database"
	"github.com/goplaceapp/goplace-guest/internal/services/guest/application"
	"gorm.io/gorm"
)

type GuestServiceServer struct {
	GuestService *application.GuestService
}

func NewGuestServiceServer(db *gorm.DB, tenantDBConnections map[string]*gorm.DB) *GuestServiceServer {
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

	return &GuestServiceServer{
		GuestService: application.NewGuestService(),
	}
}
