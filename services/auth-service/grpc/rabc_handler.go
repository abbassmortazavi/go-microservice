package grpc

import (
	rbacpb "abbassmortazavi/go-microservice/pkg/proto/rbac"
	"abbassmortazavi/go-microservice/services/auth-service/service"
	"context"
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
