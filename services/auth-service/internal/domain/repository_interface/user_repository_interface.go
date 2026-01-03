package repository_interface

import (
	"abbassmortazavi/go-microservice/services/auth-service/internal/domain/entity"
	"context"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id int64) (*entity.User, error)
}
