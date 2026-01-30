package permission

import (
	permissionpb "abbassmortazavi/go-microservice/pkg/proto/permission"
)

type DeletePermissionReq struct {
	Id int64 `json:"id" form:"id" query:"id"`
}

func (d *DeletePermissionReq) ToProto() *permissionpb.DeletePermissionRequest {
	return &permissionpb.DeletePermissionRequest{
		Id: d.Id,
	}
}
