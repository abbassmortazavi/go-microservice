package response

import (
	authpb "abbassmortazavi/go-microservice/pkg/proto/abbassmortazavi/go-microservice/auth"
	"abbassmortazavi/go-microservice/services/auth-service/entity"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type TokenResponseResult struct {
	Tokens TokenResponse `json:"tokens"`
	User   entity.User   `json:"user"`
}

func (l *TokenResponseResult) ToProto() *authpb.LoginResponse {
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
			CreatedAt: l.User.CreatedAt.Unix(),
		},
	}
}
