package permission

import (
	permissionpb "abbassmortazavi/go-microservice/pkg/proto/abbassmortazavi/go-microservice/permission"
)

type ListPermissionReq struct {
	Page    int64  `json:"page"`
	PerPage int64  `json:"per_page"`
	Search  string `json:"search"`
	SortBy  string `json:"sort_by"`
	OrderBy string `json:"order_by"`
}

func (l *ListPermissionReq) ToProto() *permissionpb.ListPermissionsRequest {
	return &permissionpb.ListPermissionsRequest{
		Page:    l.Page,
		PerPage: l.PerPage,
		Search:  l.Search,
		SortBy:  l.SortBy,
		OrderBy: l.OrderBy,
	}
}
