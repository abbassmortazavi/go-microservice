package auth

import authpb "abbassmortazavi/go-microservice/pkg/proto/abbassmortazavi/go-microservice/auth"

type GetUserRequest struct {
	ID int64 `json:"id"`
}

func (g *GetUserRequest) ToProto() *authpb.GetUserRequest {
	return &authpb.GetUserRequest{
		Id: g.ID,
	}
}
