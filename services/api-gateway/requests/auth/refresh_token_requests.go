package auth

import authpb "abbassmortazavi/go-microservice/pkg/proto/abbassmortazavi/go-microservice/auth"

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (l *RefreshTokenRequest) ToProto() *authpb.GetRefreshTokenRequest {
	return &authpb.GetRefreshTokenRequest{
		RefreshToken: l.RefreshToken,
	}
}
