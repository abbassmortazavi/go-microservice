package response

import (
	authpb "abbassmortazavi/go-microservice/pkg/proto/auth"
	"abbassmortazavi/go-microservice/services/auth-service/internal/domain/entity"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type LoginResponse struct {
	Tokens TokenResponse `json:"tokens"`
	User   entity.User   `json:"user"`
}

func (l *LoginResponse) ToProto() *authpb.LoginResponse {
	return &authpb.LoginResponse{
		Tokens: &authpb.Token{
			AccessToken:  l.Tokens.AccessToken,
			RefreshToken: l.Tokens.RefreshToken,
			ExpiredAt:    l.Tokens.ExpiresAt,
		},
		User: &authpb.User{
			Id:        l.User.ID,
			Name:      l.User.Name,
			Email:     l.User.Email,
			Role:      l.User.Role,
			CreatedAt: l.User.CreatedAt.Unix(),
		},
	}
}
