package grpc_clients

import (
	postpb "abbassmortazavi/go-microservice/pkg/proto/abbassmortazavi/go-microservice/post"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PostClientService struct {
	Client postpb.PostServiceClient
	Conn   *grpc.ClientConn
}

func NewPostClientService() (*PostClientService, error) {
	postServiceUrl := os.Getenv("POST_SERVICE_URL")
	if postServiceUrl == "" {
		postServiceUrl = "post-service:9093"
	}
	conn, err := grpc.NewClient(postServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := postpb.NewPostServiceClient(conn)
	return &PostClientService{
		Client: client,
		Conn:   conn,
	}, nil
}
