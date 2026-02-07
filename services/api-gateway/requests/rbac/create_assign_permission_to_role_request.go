package rbac

import (
	rbac "abbassmortazavi/go-microservice/pkg/proto/rbac"
)

type CreateAssignPermissionToRoleReq struct {
	RoleID       int64 `json:"role_id"`
	PermissionID int64 `json:"permission_id"`
}

func (p *CreateAssignPermissionToRoleReq) ToProto() *rbac.AssignPermissionToRoleRequest {
	return &rbac.AssignPermissionToRoleRequest{
		RoleID:       p.RoleID,
		PermissionID: p.PermissionID,
	}
}
