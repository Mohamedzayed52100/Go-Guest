package repository

import (
	"context"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-guest/database"
	"github.com/goplaceapp/goplace-guest/internal/clients/userclient"
	common "github.com/goplaceapp/goplace-guest/internal/services/common/infrastructure/repository"
	"gorm.io/gorm"
)

type ReservationLogRepository struct {
	SharedDbConnection  *gorm.DB
	TenantDBConnections map[string]*gorm.DB
	userClient          *userclient.UserClient
	CommonRepository    *common.CommonRepository
}

func NewReservationLogRepository() *ReservationLogRepository {
	postgresService := database.SharedPostgresService

	return &ReservationLogRepository{
		SharedDbConnection:  postgresService.Db,
		TenantDBConnections: postgresService.TenantDbConnections,
		CommonRepository:    common.NewCommonRepository(),
		userClient:          userclient.NewUserClient(),
	}
}

func (r *ReservationLogRepository) GetSharedDB() *gorm.DB {
	return r.SharedDbConnection
}

func (r *ReservationLogRepository) GetTenantDBConnection(ctx context.Context) *gorm.DB {
	tenantDBName := ctx.Value(meta.TenantDBNameContextKey.String()).(string)
	if db, ok := r.TenantDBConnections[tenantDBName]; ok {
		return db
	}

	return nil
}
