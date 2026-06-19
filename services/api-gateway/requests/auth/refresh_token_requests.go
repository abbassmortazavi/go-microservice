package auth

import authpb "abbassmortazavi/go-microservice/pkg/proto/abbassmortazavi/go-microservice/auth"

type RefreshTokenRequest struct {
	Token string `json:"token"`
}

func (l *RefreshTokenRequest) ToProto() *authpb.GetRefreshTokenRequest {
	return &authpb.GetRefreshTokenRequest{
		Token: l.Token,
	}
}
