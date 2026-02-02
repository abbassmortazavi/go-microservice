package role

import (
	rolepb "abbassmortazavi/go-microservice/pkg/proto/role"
)

type CreateRoleReq struct {
	Name string `json:"name"`
}

func (p *CreateRoleReq) ToProto() *rolepb.CreateRoleRequest {
	return &rolepb.CreateRoleRequest{
		Name: p.Name,
	}
}
