package repository

import (
	"context"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-guest/database"
	"github.com/goplaceapp/goplace-guest/internal/clients/userclient"
	commonRepo "github.com/goplaceapp/goplace-guest/internal/services/common/infrastructure/repository"
	"gorm.io/gorm"
)

type ReservationFeedbackRepository struct {
	SharedDbConnection  *gorm.DB
	TenantDBConnections map[string]*gorm.DB
	userClient          *userclient.UserClient
	CommonRepo          *commonRepo.CommonRepository
}

func NewReservationFeedbackRepository() *ReservationFeedbackRepository {
	postgresService := database.SharedPostgresService

	return &ReservationFeedbackRepository{
		SharedDbConnection:  postgresService.Db,
		TenantDBConnections: postgresService.TenantDbConnections,
		userClient:          userclient.NewUserClient(),
		CommonRepo:          commonRepo.NewCommonRepository(),
	}
}

func (r *ReservationFeedbackRepository) GetSharedDB() *gorm.DB {
	return r.SharedDbConnection
}

func (r *ReservationFeedbackRepository) GetTenantDBConnection(ctx context.Context) *gorm.DB {
	tenantDBName := ctx.Value(meta.TenantDBNameContextKey.String()).(string)
	if db, ok := r.TenantDBConnections[tenantDBName]; ok {
		return db
	}

	return nil
}
