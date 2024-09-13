package roleclient

import (
	"context"
	"github.com/goplaceapp/goplace-guest/database"
	"os"

	roleclient "github.com/goplaceapp/goplace-user/pkg/roleservice/adapters/grpc"
	"google.golang.org/grpc"
)

type RoleClient struct {
	Client *roleclient.RoleServiceServer
	Conn   *grpc.ClientConn
	Ctx    context.Context
}

func NewRoleClient(ctx context.Context) *RoleClient {
	conn, err := grpc.Dial(os.Getenv("USER_SERVICE_ADDRESS"), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	sharedPostgresService := database.SharedPostgresService

	return &RoleClient{
		Client: roleclient.NewRoleService(sharedPostgresService.Db, sharedPostgresService.TenantDbConnections),
		Conn:   conn,
		Ctx:    ctx,
	}
}

func (c *RoleClient) Close() {
	c.Conn.Close()
}
