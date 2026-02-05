package grpc

import (
	rbacpb "abbassmortazavi/go-microservice/pkg/proto/rbac"
	"abbassmortazavi/go-microservice/services/auth-service/service"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
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
func (r *RabcHandler) AssignPermissionToRole(ctx context.Context, req *rbacpb.AssignPermissionToRoleRequest) (*emptypb.Empty, error) {
	err := r.rbacService.AssignPermissionToRole(ctx, req.RoleID, req.PermissionID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
