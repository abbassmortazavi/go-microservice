package grpc

import (
	rbacpb "abbassmortazavi/go-microservice/pkg/proto/rbac"
	"abbassmortazavi/go-microservice/services/auth-service/service"
)

type RabcHandler struct {
	rbacpb.UnimplementedRBACServiceServer
	rbacService *service.RBACService
}

func NewRabcHandler(rb *service.RBACService) *RabcHandler {
	return &RabcHandler{
		rbacService: rb,
	}
}
