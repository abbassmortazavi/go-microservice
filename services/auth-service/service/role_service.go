package service

import (
	"abbassmortazavi/go-microservice/services/auth-service/entity"
	"abbassmortazavi/go-microservice/services/auth-service/interfaces/repository_interface"
	"context"
	"errors"
	"log"
)

type RoleService struct {
	roleRepo repository_interface.RoleRepositoryInterface
}

func NewRoleService(roleRepo repository_interface.RoleRepositoryInterface) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
	}
}
func (r *RoleService) Create(ctx context.Context, name string) (*entity.Role, error) {
	log.Println("injaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	res, _ := r.roleRepo.FindByName(ctx, name)
	if res != nil && res.Name != "" {
		return nil, errors.New(res.Name + " already exists")
	}
	role := entity.Role{
		Name: name,
	}
	roleRes, err := r.roleRepo.Save(ctx, role)
	if err != nil {
		return nil, err
	}
	res = &entity.Role{
		Name: name,
		ID:   roleRes.ID,
	}

	return res, nil
}

func (r *RoleService) Delete(ctx context.Context, permissionID int64) error {
	_, err := r.roleRepo.FindById(ctx, permissionID)
	if err != nil {
		return errors.New("role not exists")
	}
	err = r.roleRepo.Delete(ctx, permissionID)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoleService) Update(ctx context.Context, permissionID int64, name string) (entity.Role, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RoleService) Lists(ctx context.Context, page, perPage int64, orderBy, sortBy, search string) ([]entity.Role, *entity.PaginationMeta, error) {
	roles, paginationData, err := r.roleRepo.Lists(ctx, page, perPage, orderBy, sortBy, search)
	if err != nil {
		return nil, &entity.PaginationMeta{}, err
	}
	return roles, &paginationData, nil
}

func (r *RoleService) Get(ctx context.Context, permissionID int64) (*entity.Role, error) {
	role, err := r.roleRepo.FindById(ctx, permissionID)
	if err != nil {
		return nil, errors.New("role not exists")
	}
	return role, nil
}
