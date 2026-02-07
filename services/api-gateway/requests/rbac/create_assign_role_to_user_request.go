package rbac

import (
	rbac "abbassmortazavi/go-microservice/pkg/proto/rbac"
)

type CreateAssignRoleToUserReq struct {
	RoleID int64 `json:"role_id"`
	UserID int64 `json:"user_id"`
}

func (p *CreateAssignRoleToUserReq) ToProto() *rbac.AssignRoleToUserRequest {
	return &rbac.AssignRoleToUserRequest{
		RoleID: p.RoleID,
		UserID: p.UserID,
	}
}
