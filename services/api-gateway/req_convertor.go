package main

import authpb "abbassmortazavi/go-microservice/pkg/proto/auth"

type RegisterReq struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (r *RegisterReq) ToProto() *authpb.RegisterRequest {
	return &authpb.RegisterRequest{
		Email:    r.Email,
		Name:     r.Name,
		Password: r.Password,
	}
}
