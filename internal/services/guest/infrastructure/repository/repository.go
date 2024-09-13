package repository

import (
	"context"
	"github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-guest/database"
	"github.com/goplaceapp/goplace-guest/internal/clients/roleclient"
	"github.com/goplaceapp/goplace-guest/internal/clients/userclient"
	commonRepository "github.com/goplaceapp/goplace-guest/internal/services/common/infrastructure/repository"
	reservationFeedbackRepository "github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback/infrastructure/repository"

	"gorm.io/gorm"
)

type ReservationService interface {
	GetReservationByID(ctx context.Context, reservationId int32) (*domain.Reservation, error)
}

type GuestRepository struct {
	SharedDbConnection            *gorm.DB
	TenantDBConnections           map[string]*gorm.DB
	UserClient                    *userclient.UserClient
	RoleClient                    *roleclient.RoleClient
	ReservationFeedbackRepository *reservationFeedbackRepository.ReservationFeedbackRepository
	CommonRepository              *commonRepository.CommonRepository
	ReservationService
}

func NewGuestRepository() *GuestRepository {
	postgresService := database.SharedPostgresService

	return &GuestRepository{
		SharedDbConnection:            postgresService.Db,
		TenantDBConnections:           postgresService.TenantDbConnections,
		UserClient:                    userclient.NewUserClient(),
		RoleClient:                    roleclient.NewRoleClient(context.Background()),
		ReservationFeedbackRepository: reservationFeedbackRepository.NewReservationFeedbackRepository(),
		CommonRepository:              commonRepository.NewCommonRepository(),
	}
}

func (r *GuestRepository) GetSharedDB() *gorm.DB {
	return r.SharedDbConnection
}

func (r *GuestRepository) GetTenantDBConnection(ctx context.Context) *gorm.DB {
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
