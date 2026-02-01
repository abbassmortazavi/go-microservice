package grpc

import (
	rolepb "abbassmortazavi/go-microservice/pkg/proto/role"
	"abbassmortazavi/go-microservice/services/auth-service/service"
	"context"
	"log"

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

func (p *RoleHandler) Create(ctx context.Context, req *rolepb.CreateRoleRequest) (*rolepb.CreateRoleResponse, error) {
	res, err := p.roleService.Create(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	role := &rolepb.Role{
		Id:   int64(res.ID),
		Name: res.Name,
	}
	return &rolepb.CreateRoleResponse{
		Role: role,
	}, nil
}
func (p *RoleHandler) Delete(ctx context.Context, req *rolepb.DeleteRoleRequest) (*empty.Empty, error) {
	log.Println("delete role handler", req.Id)
	err := p.roleService.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
func (p *RoleHandler) Update(ctx context.Context, request *rolepb.UpdateRoleRequest) (*rolepb.UpdateRoleResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *RoleHandler) Lists(ctx context.Context, req *rolepb.ListRolesRequest) (*rolepb.ListRoleResponse, error) {
	roles, paginationData, err := p.roleService.Lists(ctx, req.Page, req.PerPage, req.OrderBy, req.SortBy, req.Search)
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

func (p *RoleHandler) Get(ctx context.Context, req *rolepb.GetRoleDetailRequest) (*rolepb.GetRoleDetailResponse, error) {
	res, err := p.roleService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	role := &rolepb.Role{
		Id:   int64(res.ID),
		Name: res.Name,
	}
	return &rolepb.GetRoleDetailResponse{
		Role: role,
	}, nil

}
