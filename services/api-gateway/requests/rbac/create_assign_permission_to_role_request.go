package rbac

import (
	rbacpb "abbassmortazavi/go-microservice/pkg/proto/abbassmortazavi/go-microservice/rbac"
)

type CreateAssignPermissionToRoleReq struct {
	RoleID       int64 `json:"role_id"`
	PermissionID int64 `json:"permission_id"`
}

func (p *CreateAssignPermissionToRoleReq) ToProto() *rbacpb.AssignPermissionToRoleRequest {
	return &rbacpb.AssignPermissionToRoleRequest{
		RoleID:       p.RoleID,
		PermissionID: p.PermissionID,
	}
}
