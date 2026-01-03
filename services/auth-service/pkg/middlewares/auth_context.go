package middlewares

import (
	"abbassmortazavi/go-microservice/services/auth-service/service"
	"context"
	"errors"
	"log"
)

var ErrUnauthorized = errors.New("unauthorized")

func User(ctx context.Context) (*service.Claims, error) {
	user, ok := ctx.Value(UserContextKey).(*service.Claims)
	log.Printf("User from: %+v\n", user)
	log.Printf("User from context: %+v\n", ctx.Value(UserContextKey))
	log.Println("ok:", ok)
	if !ok || user == nil {
		return nil, ErrUnauthorized
	}
	return user, nil
}
