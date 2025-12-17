package usecase

import (
	"abbassmortazavi/go-microservice/services/auth-service/internal/domain/entity"
	"abbassmortazavi/go-microservice/services/auth-service/internal/domain/repository"
	"context"
	"time"
)

type AuthUseCase struct {
	userRepo repository.UserRepository
}

func NewAuthUseCase(userRepo repository.UserRepository) *AuthUseCase {
	return &AuthUseCase{userRepo: userRepo}
}

func (auc *AuthUseCase) Register(ctx context.Context, email, hashPassword, name string) error {
	user := &entity.User{
		Email:     email,
		Name:      name,
		Password:  hashPassword,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return auc.userRepo.Create(ctx, user)
}
