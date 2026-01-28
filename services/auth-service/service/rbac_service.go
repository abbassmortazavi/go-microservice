package service

import (
	"abbassmortazavi/go-microservice/services/auth-service/entity"
	"abbassmortazavi/go-microservice/services/auth-service/interfaces/repository_interface"
	"context"
	"errors"
)

type RBACService struct {
	userRepo       repository_interface.UserRepositoryInterface
	roleRepo       repository_interface.RoleRepositoryInterface
	permissionRepo repository_interface.PermissionRepositoryInterface
	rbacRepo       repository_interface.RBACRepositoryInterface
}

func NewRBACService(userRepo repository_interface.UserRepositoryInterface, roleRepo repository_interface.RoleRepositoryInterface, permissionRepo repository_interface.PermissionRepositoryInterface, rbacRepo repository_interface.RBACRepositoryInterface) *RBACService {
	return &RBACService{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		rbacRepo:       rbacRepo,
	}
}
func (r *RBACService) CreatePermission(ctx context.Context, name string) (*entity.Permission, error) {
	res, err := r.permissionRepo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if res.Name != "" {
		return nil, errors.New(res.Name + " already exists")
	}
	permission := entity.Permission{
		Name: name,
	}
	p, err := r.permissionRepo.Save(ctx, permission)
	if err != nil {
		return nil, err
	}
	res = entity.Permission{
		Name: name,
		ID:   p.ID,
	}

	return &res, nil
}

func (r *RBACService) CreateRole(ctx context.Context, name string) (*entity.Role, error) {
	res, err := r.roleRepo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if res.Name != "" {
		return nil, errors.New("role already exists")
	}
	role := entity.Role{
		Name: name,
	}
	data, err := r.roleRepo.Save(ctx, &role)
	if err != nil {
		return nil, err
	}

	return data, nil
}
func (r *RBACService) AssignPermissionToRole(ctx context.Context, permissionID, roleID int) error {
	role, err := r.roleRepo.FindById(ctx, int64(roleID))
	if err != nil {
		return err
	}
	permission, err := r.permissionRepo.FindByID(ctx, permissionID)
	if err != nil {
		return err
	}
	return r.rbacRepo.AssignPermissionToRole(ctx, int64(role.ID), int64(permission.ID))
}

func (r *RBACService) AssignRoleToUser(ctx context.Context, roleID, userID int) error {
	_, err := r.roleRepo.FindById(ctx, int64(roleID))
	if err != nil {
		return err
	}
	_, err = r.userRepo.FindByID(ctx, userID)
	return r.rbacRepo.AssignRoleToUser(ctx, int64(userID), int64(roleID))
}

func (r *RBACService) CheckUserPermission(ctx context.Context, permission string, userID int64) (bool, error) {
	permissions, err := r.rbacRepo.GetPermissionsByUserID(ctx, userID)
	if err != nil {
		return false, nil
	}
	for _, p := range permissions {
		if p == permission {
			return true, nil
		}
	}
	return false, nil
}
