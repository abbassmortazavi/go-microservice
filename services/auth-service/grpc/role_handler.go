package grpc

import (
	rolepb "abbassmortazavi/go-microservice/pkg/proto/abbassmortazavi/go-microservice/role"
	"abbassmortazavi/go-microservice/services/auth-service/service"
	"context"

	"github.com/golang/protobuf/ptypes/empty"
)

type RoleHandler struct {
	rolepb.UnimplementedRoleServiceServer
	roleService *service.RoleService
}

func NewRoleHandler(p *service.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: p,
	}
}

func (r *RoleHandler) Create(ctx context.Context, req *rolepb.CreateRoleRequest) (*rolepb.CreateRoleResponse, error) {
	res, err := r.roleService.Create(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	role := &rolepb.Role{
		Id:   res.ID,
		Name: res.Name,
	}
	return &rolepb.CreateRoleResponse{
		Role: role,
	}, nil
}
func (r *RoleHandler) Delete(ctx context.Context, req *rolepb.DeleteRoleRequest) (*empty.Empty, error) {
	err := r.roleService.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
func (r *RoleHandler) Update(ctx context.Context, request *rolepb.UpdateRoleRequest) (*rolepb.UpdateRoleResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RoleHandler) Lists(ctx context.Context, req *rolepb.ListRolesRequest) (*rolepb.ListRoleResponse, error) {
	roles, paginationData, err := r.roleService.Lists(ctx, req.Page, req.PerPage, req.OrderBy, req.SortBy, req.Search)
	if err != nil {
		return nil, err
	}

	pbroles := make([]*rolepb.Role, len(roles))
	for i, role := range roles {
		pbroles[i] = &rolepb.Role{
			Id:   int64(role.ID),
			Name: role.Name,
		}
	}

	return &rolepb.ListRoleResponse{
		Roles: pbroles,
		Meta: &rolepb.Meta{
			CurrentPage: paginationData.Page,
			PerPage:     paginationData.PerPage,
			TotalPages:  paginationData.Total,
			TotalItems:  paginationData.TotalItems,
			HasNext:     paginationData.HasNextPage,
			HasPrevious: paginationData.HasPrevPage,
		},
	}, nil
}

func (r *RoleHandler) Get(ctx context.Context, req *rolepb.GetRoleDetailRequest) (*rolepb.GetRoleDetailResponse, error) {
	res, err := r.roleService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	var permissions []*rolepb.Permission
	for _, permission := range res.Permissions {
		permissions = append(permissions, &rolepb.Permission{
			Id:   int64(permission.ID),
			Name: permission.Name,
		})
	}
	result := rolepb.GetRoleDetailResponse{
		Role: &rolepb.Role{
			Id:          int64(res.ID),
			Name:        res.Name,
			Permissions: permissions,
		},
	}
	return &result, nil
}
