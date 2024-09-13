package repository

import (
	"context"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-guest/database"
	itemclient "github.com/goplaceapp/goplace-guest/internal/clients/itemclient"
	"github.com/goplaceapp/goplace-guest/internal/clients/userclient"
	commonRepo "github.com/goplaceapp/goplace-guest/internal/services/common/infrastructure/repository"
	reservationRepository "github.com/goplaceapp/goplace-guest/internal/services/reservation/infrastructure/repository"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	SharedDbConnection    *gorm.DB
	TenantDBConnections   map[string]*gorm.DB
	itemClient            *itemclient.RestaurantItemClient
	userClient            *userclient.UserClient
	commonRepository      *commonRepo.CommonRepository
	reservationRepository *reservationRepository.ReservationRepository
}

func NewPaymentRepository() *PaymentRepository {
	postgresService := database.SharedPostgresService

	return &PaymentRepository{
		SharedDbConnection:    postgresService.Db,
		TenantDBConnections:   postgresService.TenantDbConnections,
		itemClient:            itemclient.NewRestaurantItemClient(context.Background()),
		userClient:            userclient.NewUserClient(),
		commonRepository:      commonRepo.NewCommonRepository(),
		reservationRepository: reservationRepository.NewReservationRepository(),
	}
}

func (r *PaymentRepository) GetSharedDB() *gorm.DB {
	return r.SharedDbConnection
}

func (r *PaymentRepository) GetTenantDBConnection(ctx context.Context) *gorm.DB {
	tenantDBName := ctx.Value(meta.TenantDBNameContextKey.String()).(string)
	if db, ok := r.TenantDBConnections[tenantDBName]; ok {
		return db
	}

	return nil
}
