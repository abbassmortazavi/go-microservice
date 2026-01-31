package repository_interface

import (
	"abbassmortazavi/go-microservice/services/auth-service/entity"
	"context"
)

type PermissionRepositoryInterface interface {
	Create(ctx context.Context, permission entity.Permission) (*entity.Permission, error)
	FindByID(ctx context.Context, permissionId int64) (entity.Permission, error)
	FindByName(ctx context.Context, name string) (entity.Permission, error)
	Lists(ctx context.Context, page, perPage int64, orderBy, sortBy, search string) ([]entity.Permission, entity.PaginationMeta, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, name string) (entity.Permission, error)
}
