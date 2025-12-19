package service

import (
	"abbassmortazavi/go-microservice/services/auth-service/internal/domain/entity"
	"abbassmortazavi/go-microservice/services/auth-service/internal/domain/repository"
	"abbassmortazavi/go-microservice/services/auth-service/internal/infrastructure/security"
	"context"
	"errors"
	"strconv"
)

type AuthService struct {
	userRepo     repository.UserRepositoryInterface
	hasher       security.PasswordHasher
	TokenService TokenService
}

func NewAuthService(repo repository.UserRepositoryInterface, hasher security.PasswordHasher, tokenService TokenService) *AuthService {
	return &AuthService{
		userRepo:     repo,
		hasher:       hasher,
		TokenService: tokenService,
	}
}
func (a *AuthService) Register(ctx context.Context, email, password, name string) error {
	hashed, err := a.hasher.Hash(password)
	if err != nil {
		return err
	}
	user := entity.User{
		Email:    email,
		Password: hashed,
		Name:     name,
		Role:     "user",
	}
	err = a.userRepo.Create(ctx, &user)
	if err != nil {
		return err
	}
	return nil
}
func (a *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}
	if err := a.hasher.Compare(user.Password, password); err == false {
		return "", "", errors.New("password incorrect")
	}
	access, err := a.TokenService.GenerateAccessToken(strconv.FormatInt(user.ID, 10), user.Role)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := a.TokenService.GenerateRefreshToken(strconv.FormatInt(user.ID, 10))
	if err != nil {
		return "", "", err
	}
	return access, refreshToken, nil
}
