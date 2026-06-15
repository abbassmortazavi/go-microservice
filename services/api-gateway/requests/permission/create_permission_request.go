package permission

import (
	permissionpb "abbassmortazavi/go-microservice/pkg/proto/abbassmortazavi/go-microservice/permission"
)

type CreatePermissionReq struct {
	Name string `json:"name"`
}

func (p *CreatePermissionReq) ToProto() *permissionpb.CreatePermissionRequest {
	return &permissionpb.CreatePermissionRequest{
		Name: p.Name,
	}
}
