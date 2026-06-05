package rbac

import (
	rbac "abbassmortazavi/go-microservice/pkg/proto/rbac"
)

type CheckUserHasRoleReq struct {
	UserID   int64  `json:"user_id"`
	RoleName string `json:"role_name"`
}

func (p *CheckUserHasRoleReq) ToProto() *rbac.CheckUserRoleRequest {
	return &rbac.CheckUserRoleRequest{
		RoleName: p.RoleName,
		UserID:   p.UserID,
	}
}
