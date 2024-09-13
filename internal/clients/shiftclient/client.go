package shiftclient

import (
	"context"
	"os"

	"github.com/goplaceapp/goplace-guest/database"
	shiftclient "github.com/goplaceapp/goplace-settings/pkg/shiftservice/adapters/grpc"
	"google.golang.org/grpc"
)

type ShiftClient struct {
	Client *shiftclient.ShiftServiceServer
	Conn   *grpc.ClientConn
	Ctx    context.Context
}

func NewShiftClient(ctx context.Context) *ShiftClient {
	conn, err := grpc.Dial(os.Getenv("SETTINGS_SERVICE_ADDRESS"), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	postgresConn := database.SharedPostgresService

	return &ShiftClient{
		Client: shiftclient.NewShiftService(postgresConn.Db, postgresConn.TenantDbConnections),
		Conn:   conn,
		Ctx:    ctx,
	}
}

func (c *ShiftClient) Close() {
	c.Conn.Close()
}
