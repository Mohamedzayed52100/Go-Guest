package repository

import (
	"context"
	"github.com/goplaceapp/goplace-guest/internal/clients/seatingareaclient"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-guest/database"
	"github.com/goplaceapp/goplace-guest/internal/clients/shiftclient"
	commonRepo "github.com/goplaceapp/goplace-guest/internal/services/common/infrastructure/repository"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation/infrastructure/repository"

	"gorm.io/gorm"
)

type ReservationWidgetRepository struct {
	SharedDbConnection    *gorm.DB
	TenantDBConnections   map[string]*gorm.DB
	reservationRepository *repository.ReservationRepository
	commonRepository      *commonRepo.CommonRepository
	shiftClient           *shiftclient.ShiftClient
	seatingAreaClient     *seatingareaclient.SeatingAreaClient
}

func NewReservationWidgetRepository() *ReservationWidgetRepository {
	postgresService := database.SharedPostgresService

	return &ReservationWidgetRepository{
		SharedDbConnection:    postgresService.Db,
		TenantDBConnections:   postgresService.TenantDbConnections,
		reservationRepository: repository.NewReservationRepository(),
		commonRepository:      commonRepo.NewCommonRepository(),
		shiftClient:           shiftclient.NewShiftClient(context.Background()),
		seatingAreaClient:     seatingareaclient.NewSeatingAreaClient(context.Background()),
	}
}

func (r *ReservationWidgetRepository) GetSharedDB() *gorm.DB {
	return r.SharedDbConnection
}

func (r *ReservationWidgetRepository) GetTenantDBConnection(ctx context.Context) *gorm.DB {
	tenantDBName := ctx.Value(meta.TenantDBNameContextKey.String()).(string)
	if db, ok := r.TenantDBConnections[tenantDBName]; ok {
		return db
	}

	return nil
}
