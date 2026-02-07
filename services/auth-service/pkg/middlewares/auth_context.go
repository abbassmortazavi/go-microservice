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
	log.Println("ok1111111111111111111111111111111111:", ok)
	if !ok || user == nil {
		return nil, ErrUnauthorized
	}

	log.Println("hiiiiiiiiiiiiiiii:", ok)
	return user, nil
}
