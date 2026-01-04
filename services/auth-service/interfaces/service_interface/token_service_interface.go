package service_interface

import (
	"abbassmortazavi/go-microservice/services/auth-service/entity"
	"abbassmortazavi/go-microservice/services/auth-service/pkg/response"
	"abbassmortazavi/go-microservice/services/auth-service/service"
	"context"
)

type TokenServiceInterface interface {
	GenerateToken(userID int, name string) (response.TokenResponse, error)
	RefreshAccessToken(refreshToken string) (response.TokenResponse, error)
	ValidateToken(token string) (*service.Claims, error)
	FindByUserId(ctx context.Context, userId int) (*entity.Token, error)
}
