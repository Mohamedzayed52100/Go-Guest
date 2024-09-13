package repository

import (
	"context"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-guest/database"
	"github.com/goplaceapp/goplace-guest/internal/clients/userclient"
	"gorm.io/gorm"
)

type ReservationSpecialOccasionRepository struct {
	SharedDbConnection  *gorm.DB
	TenantDBConnections map[string]*gorm.DB
	userClient          *userclient.UserClient
}

func NewReservationSpecialOccasionRepository() *ReservationSpecialOccasionRepository {
	postgresService := database.SharedPostgresService

	return &ReservationSpecialOccasionRepository{
		SharedDbConnection:  postgresService.Db,
		TenantDBConnections: postgresService.TenantDbConnections,
		userClient:          userclient.NewUserClient(),
	}
}

func (r *ReservationSpecialOccasionRepository) GetSharedDB() *gorm.DB {
	return r.SharedDbConnection
}

func (r *ReservationSpecialOccasionRepository) GetTenantDBConnection(ctx context.Context) *gorm.DB {
	tenantDBNameValue := ctx.Value(meta.TenantDBNameContextKey.String())
	if tenantDBNameValue == nil {
		return nil
	}
	tenantDBName, ok := tenantDBNameValue.(string)
	if !ok {
		return nil
	}

	if db, ok := r.TenantDBConnections[tenantDBName]; ok {
		return db
	}

	return nil
}
