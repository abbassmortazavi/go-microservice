package service

import (
	eventpb "abbassmortazavi/go-microservice/pkg/proto/events"
	"abbassmortazavi/go-microservice/services/auth-service/internal/domain/entity"
	"abbassmortazavi/go-microservice/services/auth-service/internal/domain/repository_interface"
	"abbassmortazavi/go-microservice/services/auth-service/internal/infrastructure/messaging"
	"abbassmortazavi/go-microservice/services/auth-service/internal/infrastructure/security"
	"context"
	"errors"
	"strconv"
)

type AuthService struct {
	userRepo     repository_interface.UserRepositoryInterface
	hasher       security.PasswordHasher
	TokenService TokenService
	publisher    *messaging.Publisher
}

func NewAuthService(repo repository_interface.UserRepositoryInterface, hasher security.PasswordHasher, tokenService TokenService, publisher *messaging.Publisher) *AuthService {
	return &AuthService{
		userRepo:     repo,
		hasher:       hasher,
		TokenService: tokenService,
		publisher:    publisher,
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

	//
	event := eventpb.UserRegistered{
		UserId: strconv.FormatInt(user.ID, 10),
		Email:  user.Email,
		Role:   user.Role,
	}

	return a.publisher.Publish(ctx, "user_registered", event)
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
