package role

import (
	rolepb "abbassmortazavi/go-microservice/pkg/proto/role"
)

type DeleteRoleReq struct {
	Id int64 `json:"id" form:"id" query:"id"`
}

func (d *DeleteRoleReq) ToProto() *rolepb.DeleteRoleRequest {
	return &rolepb.DeleteRoleRequest{
		Id: d.Id,
	}
}
