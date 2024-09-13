package seatingareaclient

import (
	"context"
	"os"

	"github.com/goplaceapp/goplace-guest/database"
	seatingAreaClient "github.com/goplaceapp/goplace-settings/pkg/seatingareaservice/adapters/grpc"
	"google.golang.org/grpc"
)

type SeatingAreaClient struct {
	Client *seatingAreaClient.SeatingAreaServiceServer
	Conn   *grpc.ClientConn
	Ctx    context.Context
}

func NewSeatingAreaClient(ctx context.Context) *SeatingAreaClient {
	conn, err := grpc.Dial(os.Getenv("SETTINGS_SERVICE_ADDRESS"), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	postgresConn := database.SharedPostgresService

	return &SeatingAreaClient{
		Client: seatingAreaClient.NewSeatingAreaService(postgresConn.Db, postgresConn.TenantDbConnections),
		Conn:   conn,
		Ctx:    ctx,
	}
}

func (c *SeatingAreaClient) Close() {
	c.Conn.Close()
}
