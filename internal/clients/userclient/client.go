package userclient

import (
	"github.com/goplaceapp/goplace-guest/database"
	userclient "github.com/goplaceapp/goplace-user/pkg/userservice/adapters/grpc"
	"google.golang.org/grpc"
	"os"
)

type UserClient struct {
	Client *userclient.UserServiceServer
	Conn   *grpc.ClientConn
}

func NewUserClient() *UserClient {
	conn, err := grpc.Dial(os.Getenv("USER_SERVICE_ADDRESS"), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	sharedPostgresService := database.SharedPostgresService

	return &UserClient{
		Client: userclient.NewUserService(sharedPostgresService.Db, sharedPostgresService.TenantDbConnections),
		Conn:   conn,
	}
}

func (c *UserClient) Close() {
	c.Conn.Close()
}
