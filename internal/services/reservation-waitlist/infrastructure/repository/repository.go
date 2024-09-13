package repository

import (
	"context"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-guest/database"
	"github.com/goplaceapp/goplace-guest/internal/clients/seatingareaclient"
	"github.com/goplaceapp/goplace-guest/internal/clients/shiftclient"
	"github.com/goplaceapp/goplace-guest/internal/clients/userclient"
	commonRepo "github.com/goplaceapp/goplace-guest/internal/services/common/infrastructure/repository"
	guestRepostiory "github.com/goplaceapp/goplace-guest/internal/services/guest/infrastructure/repository"
	reservationRepository "github.com/goplaceapp/goplace-guest/internal/services/reservation/infrastructure/repository"
	"gorm.io/gorm"
)

type ReservationWaitListRepository struct {
	SharedDbConnection    *gorm.DB
	TenantDBConnections   map[string]*gorm.DB
	userClient            *userclient.UserClient
	reservationRepository *reservationRepository.ReservationRepository
	guestRepository       *guestRepostiory.GuestRepository
	seatingAreaClient     *seatingareaclient.SeatingAreaClient
	shiftClient           *shiftclient.ShiftClient
	CommonRepo            *commonRepo.CommonRepository
}

func NewReservationWaitListRepository() *ReservationWaitListRepository {
	postgresService := database.SharedPostgresService

	return &ReservationWaitListRepository{
		SharedDbConnection:    postgresService.Db,
		TenantDBConnections:   postgresService.TenantDbConnections,
		userClient:            userclient.NewUserClient(),
		reservationRepository: reservationRepository.NewReservationRepository(),
		guestRepository:       guestRepostiory.NewGuestRepository(),
		seatingAreaClient:     seatingareaclient.NewSeatingAreaClient(context.Background()),
		shiftClient:           shiftclient.NewShiftClient(context.Background()),
		CommonRepo:            commonRepo.NewCommonRepository(),
	}
}

func (r *ReservationWaitListRepository) GetSharedDB() *gorm.DB {
	return r.SharedDbConnection
}

func (r *ReservationWaitListRepository) GetTenantDBConnection(ctx context.Context) *gorm.DB {
	tenantDBName := ctx.Value(meta.TenantDBNameContextKey.String()).(string)
	if db, ok := r.TenantDBConnections[tenantDBName]; ok {
		return db
	}

	return nil
}
