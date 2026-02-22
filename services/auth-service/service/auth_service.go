package service

import (
	eventpb "abbassmortazavi/go-microservice/pkg/proto/events"
	"abbassmortazavi/go-microservice/services/auth-service/entity"
	"abbassmortazavi/go-microservice/services/auth-service/interfaces/repository_interface"
	"abbassmortazavi/go-microservice/services/auth-service/messaging"
	"abbassmortazavi/go-microservice/services/auth-service/pkg/response"
	"abbassmortazavi/go-microservice/services/auth-service/security"
	"context"
	"errors"
	"strconv"
)

type AuthServiceInterface interface {
	Register(ctx context.Context, email, password, name string) error
	Login(ctx context.Context, email, password string) (*response.LoginResponse, error)
	GetUser(ctx context.Context, id int64) (*entity.User, error)
	CheckUseHasRole(ctx context.Context, userID int64, requiredRoles []string) (bool, error)
	CheckSingleRole(ctx context.Context, userID int64, roleName string) (bool, error)
}

type AuthService struct {
	userRepo     repository_interface.UserRepositoryInterface
	hasher       security.PasswordHasher
	TokenService TokenServiceInterface
	publisher    *messaging.Publisher
}

func NewAuthService(repo repository_interface.UserRepositoryInterface, hasher security.PasswordHasher, tokenService TokenServiceInterface, publisher *messaging.Publisher) *AuthService {
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
	}
	err = a.userRepo.Create(ctx, &user)
	if err != nil {
		return err
	}

	//
	event := eventpb.UserRegistered{
		UserId: strconv.FormatInt(user.ID, 10),
		Email:  user.Email,
	}

	return a.publisher.Publish(ctx, "user_registered", event)
}
func (a *AuthService) Login(ctx context.Context, email, password string) (*response.LoginResponse, error) {
	user, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if err := a.hasher.Compare(user.Password, password); err == false {
		return nil, errors.New("password is wrong")
	}
	tokens, err := a.TokenService.GenerateToken(user.ID, user.Name)
	if err != nil {
		return nil, err
	}

	userEntity := entity.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	return &response.LoginResponse{
		Tokens: tokens,
		User:   userEntity,
	}, nil
}
func (a *AuthService) GetUser(ctx context.Context, id int64) (*entity.User, error) {
	user, err := a.userRepo.FindByID(ctx, id)

	if err != nil {
		return nil, err
	}
	/*userEntity := entity.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}*/
	return user, nil
}

func (a *AuthService) CheckUseHasRole(ctx context.Context, userID int64, requiredRoles []string) (bool, error) {
	user, err := a.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, nil
	}
	if user.ID != userID {
		return false, nil
	}
	userRoles := make(map[string]bool)
	for _, role := range user.Role {
		userRoles[role.Name] = true
	}
	for _, role := range requiredRoles {
		if !userRoles[role] {
			return false, nil
		}
	}
	return true, nil
}
func (a *AuthService) CheckSingleRole(ctx context.Context, userID int64, roleName string) (bool, error) {
	user, err := a.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, nil
	}
	if user.ID != userID {
		return false, nil
	}
	for _, role := range user.Role {
		if role.Name == roleName {
			return true, nil
		}
	}
	return false, nil
}
