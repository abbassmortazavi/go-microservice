package permission

import rbacpb "abbassmortazavi/go-microservice/pkg/proto/rbac"

type DeletePermissionReq struct {
	Id int64 `json:"id" form:"id" query:"id"`
}

func (d *DeletePermissionReq) ToProto() *rbacpb.DeletePermissionRequest {
	return &rbacpb.DeletePermissionRequest{
		Id: d.Id,
	}
}
