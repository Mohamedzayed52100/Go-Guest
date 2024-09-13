package repository

import (
	"context"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-guest/database"
	"gorm.io/gorm"
)

type ReservationFeedbackWebhookRepository struct {
	SharedDbConnection  *gorm.DB
	TenantDBConnections map[string]*gorm.DB
}

func NewReservationFeedbackWebhookRepository() *ReservationFeedbackWebhookRepository {
	postgresService := database.SharedPostgresService

	return &ReservationFeedbackWebhookRepository{
		SharedDbConnection:  postgresService.Db,
		TenantDBConnections: postgresService.TenantDbConnections,
	}
}

func (r *ReservationFeedbackWebhookRepository) GetSharedDB() *gorm.DB {
	return r.SharedDbConnection
}

func (r *ReservationFeedbackWebhookRepository) GetTenantDBConnection(ctx context.Context) *gorm.DB {
	tenantDBName := ctx.Value(meta.TenantDBNameContextKey.String()).(string)
	if db, ok := r.TenantDBConnections[tenantDBName]; ok {
		return db
	}

	return nil
}
