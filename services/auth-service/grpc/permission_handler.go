package grpc

import (
	permissionpb "abbassmortazavi/go-microservice/pkg/proto/abbassmortazavi/go-microservice/permission"
	"abbassmortazavi/go-microservice/services/auth-service/service"
	"context"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
)

type PermissionHandler struct {
	permissionpb.UnimplementedPermissionServiceServer
	permissionService *service.PermissionService
}

func NewPermissionHandler(p *service.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		permissionService: p,
	}
}

func (p *PermissionHandler) Create(ctx context.Context, req *permissionpb.CreatePermissionRequest) (*permissionpb.CreatePermissionResponse, error) {
	res, err := p.permissionService.Create(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	permission := &permissionpb.Permission{
		Id:   int64(res.ID),
		Name: res.Name,
	}
	return &permissionpb.CreatePermissionResponse{
		Permission: permission,
	}, nil
}
func (p *PermissionHandler) Delete(ctx context.Context, req *permissionpb.DeletePermissionRequest) (*empty.Empty, error) {
	log.Println("delete permission handler", req.Id)
	err := p.permissionService.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
func (p *PermissionHandler) Update(ctx context.Context, request *permissionpb.UpdatePermissionRequest) (*permissionpb.UpdatePermissionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PermissionHandler) Lists(ctx context.Context, req *permissionpb.ListPermissionsRequest) (*permissionpb.ListPermissionResponse, error) {
	permissions, paginationData, err := p.permissionService.Lists(ctx, req.Page, req.PerPage, req.OrderBy, req.SortBy, req.Search)
	if err != nil {
		return nil, err
	}

	pbpermissions := make([]*permissionpb.Permission, len(permissions))
	for i, permission := range permissions {
		pbpermissions[i] = &permissionpb.Permission{
			Id:   int64(permission.ID),
			Name: permission.Name,
		}
	}

	return &permissionpb.ListPermissionResponse{
		Permissions: pbpermissions,
		Meta: &permissionpb.Meta{
			CurrentPage: paginationData.Page,
			PerPage:     paginationData.PerPage,
			TotalPages:  paginationData.Total,
			TotalItems:  paginationData.TotalItems,
			HasNext:     paginationData.HasNextPage,
			HasPrevious: paginationData.HasPrevPage,
		},
	}, nil
}

func (p *PermissionHandler) Get(ctx context.Context, req *permissionpb.GetPermissionDetailRequest) (*permissionpb.GetPermissionDetailResponse, error) {
	res, err := p.permissionService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	permission := &permissionpb.Permission{
		Id:   int64(res.ID),
		Name: res.Name,
	}
	return &permissionpb.GetPermissionDetailResponse{
		Permission: permission,
	}, nil

}
