package common

import (
	"context"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-guest/database"
	"github.com/goplaceapp/goplace-guest/internal/clients/seatingareaclient"
	"github.com/goplaceapp/goplace-guest/internal/clients/shiftclient"
	"github.com/goplaceapp/goplace-guest/internal/clients/userclient"
	"gorm.io/gorm"
)

type CommonRepository struct {
	SharedDbConnection  *gorm.DB
	TenantDBConnections map[string]*gorm.DB
	userClient          *userclient.UserClient
	seatingAreaClient   *seatingareaclient.SeatingAreaClient
	shiftClient         *shiftclient.ShiftClient
}

func NewCommonRepository() *CommonRepository {
	postgresService := database.SharedPostgresService

	return &CommonRepository{
		SharedDbConnection:  postgresService.Db,
		TenantDBConnections: postgresService.TenantDbConnections,
		userClient:          userclient.NewUserClient(),
		seatingAreaClient:   seatingareaclient.NewSeatingAreaClient(context.Background()),
		shiftClient:         shiftclient.NewShiftClient(context.Background()),
	}
}

func (r *CommonRepository) GetSharedDB() *gorm.DB {
	return r.SharedDbConnection
}

func (r *CommonRepository) GetTenantDBConnection(ctx context.Context) *gorm.DB {
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
