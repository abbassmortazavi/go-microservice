package service

import (
	"abbassmortazavi/go-microservice/services/auth-service/entity"
	"abbassmortazavi/go-microservice/services/auth-service/interfaces/repository_interface"
	"context"
	"errors"
	"log"
)

type PermissionService struct {
	permissionRepo repository_interface.PermissionRepositoryInterface
}

func NewPermissionService(permissionRepo repository_interface.PermissionRepositoryInterface) *PermissionService {
	return &PermissionService{
		permissionRepo: permissionRepo,
	}
}
func (p *PermissionService) Create(ctx context.Context, name string) (*entity.Permission, error) {
	res, _ := p.permissionRepo.FindByName(ctx, name)
	if res.Name != "" {
		return nil, errors.New(res.Name + " already exists")
	}
	permission := entity.Permission{
		Name: name,
	}
	permissionRes, err := p.permissionRepo.Create(ctx, permission)
	if err != nil {
		return nil, err
	}
	res = entity.Permission{
		Name: name,
		ID:   permissionRes.ID,
	}

	return &res, nil
}

func (p *PermissionService) Delete(ctx context.Context, permissionID int64) error {
	_, err := p.permissionRepo.FindByID(ctx, permissionID)
	if err != nil {
		return errors.New("permission not exists")
	}
	err = p.permissionRepo.Delete(ctx, permissionID)
	if err != nil {
		return err
	}
	return nil
}

func (p *PermissionService) Update(ctx context.Context, permissionID int64, name string) (entity.Permission, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PermissionService) Lists(ctx context.Context, page, perPage int64, orderBy, sortBy, search string) ([]entity.Permission, *entity.PaginationMeta, error) {
	log.Println("permission service 3")
	permissions, paginationData, err := p.permissionRepo.Lists(ctx, page, perPage, orderBy, sortBy, search)
	if err != nil {
		return nil, &entity.PaginationMeta{}, err
	}
	log.Println("permission lists", permissions)
	log.Println("pagination data", paginationData)
	return permissions, &paginationData, nil
}
