package tableclient

import (
	"context"
	"os"

	"github.com/goplaceapp/goplace-guest/database"
	tableservice "github.com/goplaceapp/goplace-settings/pkg/tableservice/adapters/grpc"
	"google.golang.org/grpc"
)

type TableClient struct {
	Client *tableservice.TableServiceServer
	Conn   *grpc.ClientConn
	Ctx    context.Context
}

func NewTableClient(ctx context.Context) *TableClient {
	conn, err := grpc.Dial(os.Getenv("SETTINGS_SERVICE_ADDRESS"), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	postgresConn := database.SharedPostgresService

	return &TableClient{
		Client: tableservice.NewTableService(postgresConn.Db, postgresConn.TenantDbConnections),
		Conn:   conn,
		Ctx:    ctx,
	}
}
