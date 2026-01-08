package repository_interface

import (
	"abbassmortazavi/go-microservice/services/auth-service/entity"
	"context"
)

type RoleRepositoryInterface interface {
	Save(ctx context.Context, role *entity.Role) error
	FindById(ctx context.Context, roleId int64) (*entity.Role, error)
	FindByName(ctx context.Context, name string) (*entity.Role, error)
	Lists(ctx context.Context) ([]*entity.Role, error)
	Delete(ctx context.Context, roleId int64) error
}
