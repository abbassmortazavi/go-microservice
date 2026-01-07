package service

import (
	"abbassmortazavi/go-microservice/services/auth-service/entity"
	"abbassmortazavi/go-microservice/services/auth-service/repository"
	"context"
	"errors"
)

type RBACService struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
	rbacRepo       repository.RBACRepository
}

func NewRBACService(userRepo repository.UserRepository, roleRepo repository.RoleRepository, permissionRepo repository.PermissionRepository, rbacRepo repository.RBACRepository) *RBACService {
	return &RBACService{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		rbacRepo:       rbacRepo,
	}
}
func (r *RBACService) CreatePermission(ctx context.Context, name string) error {
	res, err := r.permissionRepo.FindByName(ctx, name)
	if err != nil {
		return err
	}
	if res.Name != "" {
		return errors.New("permission already exists")
	}
	permission := entity.Permission{
		Name: name,
	}
	return r.permissionRepo.Save(ctx, permission)
}
