package grpc

import (
	authpb "abbassmortazavi/go-microservice/pkg/proto/auth"
	"abbassmortazavi/go-microservice/services/auth-service/internal/usecase"
	"context"
)

type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	usecase *usecase.AuthUseCase
}

func NewAuthHandler(u *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		usecase: u,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	// password hashing later
	err := h.usecase.Register(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		return nil, err
	}

	return &authpb.RegisterResponse{Message: "registered"}, nil
}
