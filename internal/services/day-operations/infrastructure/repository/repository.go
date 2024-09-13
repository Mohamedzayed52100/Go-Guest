package repository

import (
	"context"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-guest/database"
	"github.com/goplaceapp/goplace-guest/internal/clients/userclient"
	commonRepo "github.com/goplaceapp/goplace-guest/internal/services/common/infrastructure/repository"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation/infrastructure/repository"
	"gorm.io/gorm"
)

type DayOperationsRepository struct {
	SharedDbConnection    *gorm.DB
	TenantDBConnections   map[string]*gorm.DB
	userClient            *userclient.UserClient
	CommonRepo            *commonRepo.CommonRepository
	reservationRepository *repository.ReservationRepository
}

func NewDayOperationsRepository() *DayOperationsRepository {
	postgresService := database.SharedPostgresService

	return &DayOperationsRepository{
		SharedDbConnection:    postgresService.Db,
		TenantDBConnections:   postgresService.TenantDbConnections,
		userClient:            userclient.NewUserClient(),
		CommonRepo:            commonRepo.NewCommonRepository(),
		reservationRepository: repository.NewReservationRepository(),
	}
}

func (r *DayOperationsRepository) GetSharedDB() *gorm.DB {
	return r.SharedDbConnection
}

func (r *DayOperationsRepository) GetTenantDBConnection(ctx context.Context) *gorm.DB {
	tenantDBName := ctx.Value(meta.TenantDBNameContextKey.String()).(string)
	if db, ok := r.TenantDBConnections[tenantDBName]; ok {
		return db
	}

	return nil
}
