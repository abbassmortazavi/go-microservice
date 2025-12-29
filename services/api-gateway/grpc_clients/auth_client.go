package grpc_clients

import (
	authpb "abbassmortazavi/go-microservice/pkg/proto/auth"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthServiceClient struct {
	Client authpb.AuthServiceClient
	Conn   *grpc.ClientConn
}

func NewAuthServiceClient() (*AuthServiceClient, error) {
	authServiceUrl := os.Getenv("AUTH_SERVICE_URL")
	if authServiceUrl == "" {
		authServiceUrl = "auth-service:9092"
	}
	conn, err := grpc.NewClient(authServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := authpb.NewAuthServiceClient(conn)
	return &AuthServiceClient{
		Client: client,
		Conn:   conn,
	}, nil
}
func (c *AuthServiceClient) Close() error {
	return c.Conn.Close()
}
