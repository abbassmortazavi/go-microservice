package repository_interface

import (
	"abbassmortazavi/go-microservice/services/auth-service/entity"
	"context"
)

type TokenRepositoryInterface interface {
	Create(ctx context.Context, token *entity.Token) error
	FindByToken(ctx context.Context, token string) (*entity.Token, error)
	Revoke(ctx context.Context, token string) error
	FindByUserId(ctx context.Context, userId int) (*entity.Token, error)
	Delete(ctx context.Context, id int) error
	RevokeAllUserTokens(ctx context.Context, userId int) error
}
