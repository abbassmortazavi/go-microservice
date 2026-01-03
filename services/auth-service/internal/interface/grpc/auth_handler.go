package grpc

import (
	authpb "abbassmortazavi/go-microservice/pkg/proto/auth"
	"abbassmortazavi/go-microservice/services/auth-service/internal/domain/service"
	"context"
	"log"
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
	log.Println(2)
	err := h.authService.Register(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		return nil, err
	}

	return &authpb.RegisterResponse{Message: "registered"}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	log.Println("im here")
	res, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return res.ToProto(), nil
	/*return &authpb.LoginResponse{
		Tokens: &authpb.Token{
			AccessToken:  res.Tokens.AccessToken,
			RefreshToken: res.Tokens.RefreshToken,
			ExpiredAt:    res.Tokens.ExpiresAt,
		},
		User: &authpb.User{
			Id:        res.User.ID,
			Name:      res.User.Name,
			Email:     res.User.Email,
			Role:      res.User.Role,
			CreatedAt: res.User.CreatedAt.Unix(),
		},
	}, nil*/
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

func (h *AuthHandler) GetUser(ctx context.Context, req *authpb.GetUserRequest) (*authpb.GetUserResponse, error) {
	res, err := h.authService.GetUser(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &authpb.GetUserResponse{
		User: &authpb.User{
			Id:        res.ID,
			Name:      res.Name,
			Email:     res.Email,
			Role:      res.Role,
			CreatedAt: res.CreatedAt.Unix(),
		},
	}, nil
}
