package rbac

import (
	rbacpb "abbassmortazavi/go-microservice/pkg/proto/abbassmortazavi/go-microservice/rbac"
)

type CreateAssignRoleToUserReq struct {
	RoleID int64 `json:"role_id"`
	UserID int64 `json:"user_id"`
}

func (p *CreateAssignRoleToUserReq) ToProto() *rbacpb.AssignRoleToUserRequest {
	return &rbacpb.AssignRoleToUserRequest{
		RoleID: p.RoleID,
		UserID: p.UserID,
	}
}
