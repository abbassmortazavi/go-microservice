package requests

import authpb "abbassmortazavi/go-microservice/pkg/proto/auth"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *LoginRequest) ToProto() *authpb.LoginRequest {
	return &authpb.LoginRequest{
		Email:    l.Email,
		Password: l.Password,
	}
}
