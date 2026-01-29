package grpc

import (
	rbacpb "abbassmortazavi/go-microservice/pkg/proto/rbac"
	"abbassmortazavi/go-microservice/services/auth-service/service"
	"context"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
)

type RabcHandler struct {
	rbacpb.UnimplementedRBACServiceServer
	rbacService *service.RBACService
}

func NewRabcHandler(rb *service.RBACService) *RabcHandler {
	return &RabcHandler{
		rbacService: rb,
	}
}

func (rb *RabcHandler) CreatePermission(ctx context.Context, req *rbacpb.CreatePermissionRequest) (*rbacpb.CreatePermissionResponse, error) {
	res, err := rb.rbacService.CreatePermission(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	permission := &rbacpb.Permission{
		Id:   int64(res.ID),
		Name: res.Name,
	}
	return &rbacpb.CreatePermissionResponse{
		Permission: permission,
	}, nil
}
func (rb *RabcHandler) DeletePermission(ctx context.Context, req *rbacpb.DeletePermissionRequest) (*empty.Empty, error) {
	log.Println("delete permission handler", req.Id)
	err := rb.rbacService.DeletePermission(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
