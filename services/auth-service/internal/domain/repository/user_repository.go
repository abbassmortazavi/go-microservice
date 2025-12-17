package repository

import (
	"abbassmortazavi/go-microservice/services/auth-service/internal/domain/entity"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}
