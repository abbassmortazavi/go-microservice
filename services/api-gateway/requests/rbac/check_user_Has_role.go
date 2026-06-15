package rbac

import (
	rbacpb "abbassmortazavi/go-microservice/pkg/proto/abbassmortazavi/go-microservice/rbac"
)

type CheckUserHasRoleReq struct {
	UserID   int64  `json:"user_id"`
	RoleName string `json:"role_name"`
}

func (p *CheckUserHasRoleReq) ToProto() *rbacpb.CheckUserRoleRequest {
	return &rbacpb.CheckUserRoleRequest{
		RoleName: p.RoleName,
		UserID:   p.UserID,
	}
}
