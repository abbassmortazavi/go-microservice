package rbac

import (
	rbac "abbassmortazavi/go-microservice/pkg/proto/rbac"
)

type CheckUserHasPermissionReq struct {
	UserID         int64  `json:"user_id"`
	PermissionName string `json:"permission_name"`
}

func (p *CheckUserHasPermissionReq) ToProto() *rbac.CheckUserPermissionRequest {
	return &rbac.CheckUserPermissionRequest{
		PermissionName: p.PermissionName,
		UserID:         p.UserID,
	}
}
