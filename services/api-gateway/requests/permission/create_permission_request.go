package permission

import rbacpb "abbassmortazavi/go-microservice/pkg/proto/rbac"

type CreatePermissionReq struct {
	Name string `json:"name"`
}

func (p *CreatePermissionReq) ToProto() *rbacpb.CreatePermissionRequest {
	return &rbacpb.CreatePermissionRequest{
		Name: p.Name,
	}
}
