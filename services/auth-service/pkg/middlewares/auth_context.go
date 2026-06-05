package middlewares

import (
	"abbassmortazavi/go-microservice/services/auth-service/service"
	"context"
	"errors"
)

var ErrUnauthorized = errors.New("unauthorized")

func User(ctx context.Context) (*service.Claims, error) {
	user, ok := ctx.Value(UserContextKey).(*service.Claims)
	if !ok || user == nil {
		return nil, ErrUnauthorized
	}
	return user, nil
}
