package repository

import (
	"context"

	"github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/database"
	"github.com/goplaceapp/goplace-guest/internal/clients/roleclient"
	"github.com/goplaceapp/goplace-guest/internal/clients/seatingareaclient"
	"github.com/goplaceapp/goplace-guest/internal/clients/shiftclient"
	"github.com/goplaceapp/goplace-guest/internal/clients/tableclient"
	"github.com/goplaceapp/goplace-guest/internal/clients/userclient"
	commonRepo "github.com/goplaceapp/goplace-guest/internal/services/common/infrastructure/repository"

	"gorm.io/gorm"
)

type GuestService interface {
	GetGuestByID(ctx context.Context, req *guestProto.GetGuestByIDRequest) (*domain.Guest, error)
}

type ReservationRepository struct {
	SharedDbConnection  *gorm.DB
	TenantDBConnections map[string]*gorm.DB
	userClient          *userclient.UserClient
	roleClient          *roleclient.RoleClient
	shiftClient         *shiftclient.ShiftClient
	seatingAreaClient   *seatingareaclient.SeatingAreaClient
	tableClient         *tableclient.TableClient
	GuestService
	CommonRepo *commonRepo.CommonRepository
}

func NewReservationRepository() *ReservationRepository {
	postgresService := database.SharedPostgresService

	return &ReservationRepository{
		SharedDbConnection:  postgresService.Db,
		TenantDBConnections: postgresService.TenantDbConnections,
		userClient:          userclient.NewUserClient(),
		shiftClient:         shiftclient.NewShiftClient(context.Background()),
		seatingAreaClient:   seatingareaclient.NewSeatingAreaClient(context.Background()),
		roleClient:          roleclient.NewRoleClient(context.Background()),
		tableClient:         tableclient.NewTableClient(context.Background()),
		CommonRepo:          commonRepo.NewCommonRepository(),
	}
}

func (r *ReservationRepository) GetSharedDB() *gorm.DB {
	return r.SharedDbConnection
}

func (r *ReservationRepository) GetTenantDBConnection(ctx context.Context) *gorm.DB {
	tenantDBName := ctx.Value(meta.TenantDBNameContextKey.String()).(string)
	if db, ok := r.TenantDBConnections[tenantDBName]; ok {
		return db
	}

	return nil
}
