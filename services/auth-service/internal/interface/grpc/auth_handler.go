package grpc

import (
	authpb "abbassmortazavi/go-microservice/pkg/proto/auth"
	"abbassmortazavi/go-microservice/services/auth-service/internal/domain/service"
	"context"
)

type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	authService *service.AuthService
}

func NewAuthHandler(a *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: a,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	// password hashing later
	err := h.authService.Register(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		return nil, err
	}

	return &authpb.RegisterResponse{Message: "registered"}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	access, refresh, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &authpb.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
func (h *AuthHandler) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	claims, err := h.authService.TokenService.ValidateToken(req.Token)
	if err != nil {
		return nil, err
	}
	return &authpb.ValidateTokenResponse{
		UserId: claims.Subject,
	}, nil
}
