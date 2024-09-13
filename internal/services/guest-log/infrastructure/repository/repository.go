package repository

import (
	"context"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-guest/database"
	"github.com/goplaceapp/goplace-guest/internal/clients/userclient"
	common "github.com/goplaceapp/goplace-guest/internal/services/common/infrastructure/repository"
	"gorm.io/gorm"
)

type GuestLogRepository struct {
	SharedDbConnection  *gorm.DB
	TenantDBConnections map[string]*gorm.DB
	userClient          *userclient.UserClient
	CommonRepository    *common.CommonRepository
}

func NewGuestLogRepository() *GuestLogRepository {
	postgresService := database.SharedPostgresService

	return &GuestLogRepository{
		SharedDbConnection:  postgresService.Db,
		TenantDBConnections: postgresService.TenantDbConnections,
		CommonRepository:    common.NewCommonRepository(),
		userClient:          userclient.NewUserClient(),
	}
}

func (r *GuestLogRepository) GetSharedDB() *gorm.DB {
	return r.SharedDbConnection
}

func (r *GuestLogRepository) GetTenantDBConnection(ctx context.Context) *gorm.DB {
	tenantDBName := ctx.Value(meta.TenantDBNameContextKey.String()).(string)
	if db, ok := r.TenantDBConnections[tenantDBName]; ok {
		return db
	}

	return nil
}
