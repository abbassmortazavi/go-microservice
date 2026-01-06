package repository_interface

import (
	"abbassmortazavi/go-microservice/services/auth-service/entity"
	"context"
)

type PermissionRepositoryInterface interface {
	Save(ctx context.Context, permission entity.Permission) error
	FindByID(ctx context.Context, permissionId int) (entity.Permission, error)
	FindByName(ctx context.Context, name string) (entity.Permission, error)
	Lists(ctx context.Context) ([]entity.Permission, error)
}
