package shiftclient

import (
	"context"
	"os"

	"github.com/goplaceapp/goplace-guest/database"
	itemclient "github.com/goplaceapp/goplace-settings/pkg/itemservice/adapters/grpc"
	"google.golang.org/grpc"
)

type RestaurantItemClient struct {
	Client *itemclient.RestaurantItemServiceServer
	Conn   *grpc.ClientConn
	Ctx    context.Context
}

func NewRestaurantItemClient(ctx context.Context) *RestaurantItemClient {
	conn, err := grpc.Dial(os.Getenv("SETTINGS_SERVICE_ADDRESS"), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	postgresConn := database.SharedPostgresService

	return &RestaurantItemClient{
		Client: itemclient.NewRestaurantItemService(postgresConn.Db, postgresConn.TenantDbConnections),
		Conn:   conn,
		Ctx:    ctx,
	}
}

func (c *RestaurantItemClient) Close() {
	c.Conn.Close()
}
