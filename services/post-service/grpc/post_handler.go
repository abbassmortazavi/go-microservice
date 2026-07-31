package grpc

import (
	postpb "abbassmortazavi/go-microservice/pkg/proto/abbassmortazavi/go-microservice/post"
	"context"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"
)

type PostHandler struct {
	postpb.UnimplementedPostServiceServer
}

func NewPostHandler() *PostHandler {
	return &PostHandler{}
}

func (p *PostHandler) List(ctx context.Context, empty *emptypb.Empty) (*postpb.ListResponse, error) {
	log.Println("PostHandler.List called")
	return &postpb.ListResponse{
		Message: "hiiiiiii",
	}, nil
}
