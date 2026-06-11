package grpc

import (
	rbacpb "abbassmortazavi/go-microservice/pkg/proto/abbassmortazavi/go-microservice/rbac"
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
func (r *RabcHandler) AssignRoleToUser(ctx context.Context, req *rbacpb.AssignRoleToUserRequest) (*emptypb.Empty, error) {
	err := r.rbacService.AssignRoleToUser(ctx, req.GetRoleID(), req.GetUserID())
	if err != nil {
		return nil, err
	}
	return nil, nil
}
func (r *RabcHandler) CheckUserRole(ctx context.Context, req *rbacpb.CheckUserRoleRequest) (*rbacpb.CheckUserRoleResponse, error) {
	res, err := r.rbacService.CheckUserHasRole(ctx, req.GetRoleName(), req.GetUserID())
	if err != nil {
		return &rbacpb.CheckUserRoleResponse{
			HasRole: false,
		}, err
	}

	return &rbacpb.CheckUserRoleResponse{
		HasRole: res,
	}, nil
}
