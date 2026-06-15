package grpc_clients

import (
	authpb "abbassmortazavi/go-microservice/pkg/proto/abbassmortazavi/go-microservice/auth"
	permissionpb "abbassmortazavi/go-microservice/pkg/proto/abbassmortazavi/go-microservice/permission"
	rbacpb "abbassmortazavi/go-microservice/pkg/proto/abbassmortazavi/go-microservice/rbac"
	rolepb "abbassmortazavi/go-microservice/pkg/proto/abbassmortazavi/go-microservice/role"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthServiceClient struct {
	Client     authpb.AuthServiceClient
	Conn       *grpc.ClientConn
	Rbac       rbacpb.RBACServiceClient
	Permission permissionpb.PermissionServiceClient
	Role       rolepb.RoleServiceClient
}

func NewAuthServiceClient() (*AuthServiceClient, error) {
	authServiceUrl := os.Getenv("AUTH_SERVICE_URL")
	if authServiceUrl == "" {
		authServiceUrl = "auth-service:9092"
	}
	conn, err := grpc.NewClient(authServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := authpb.NewAuthServiceClient(conn)
	rbac := rbacpb.NewRBACServiceClient(conn)
	permission := permissionpb.NewPermissionServiceClient(conn)
	role := rolepb.NewRoleServiceClient(conn)
	return &AuthServiceClient{
		Client:     client,
		Conn:       conn,
		Rbac:       rbac,
		Permission: permission,
		Role:       role,
	}, nil
}
func (c *AuthServiceClient) Close() error {
	return c.Conn.Close()
}
