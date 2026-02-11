package implement

import (
	global "abbassmortazavi/go-microservice/pkg/config"
	"abbassmortazavi/go-microservice/pkg/database"
	"abbassmortazavi/go-microservice/pkg/message"
	authpb "abbassmortazavi/go-microservice/pkg/proto/auth"
	permissionpb "abbassmortazavi/go-microservice/pkg/proto/permission"
	rbacpb "abbassmortazavi/go-microservice/pkg/proto/rbac"
	rolepb "abbassmortazavi/go-microservice/pkg/proto/role"
	"abbassmortazavi/go-microservice/services/auth-service/grpc"
	"abbassmortazavi/go-microservice/services/auth-service/pkg/middlewares"
	"abbassmortazavi/go-microservice/services/auth-service/repository"
	"abbassmortazavi/go-microservice/services/auth-service/security"
	service2 "abbassmortazavi/go-microservice/services/auth-service/service"

	grpc2 "google.golang.org/grpc"
)

var Server *grpc2.Server

func Init() {
	gcfg := global.Load()
	database.Connect()

	// ---- RabbitMQ ----
	message.Init()
	publisher := message.GetPublisher()

	userRepo := repository.NewUserRepository(database.DB)
	tokenRepo := repository.NewTokenRepository(database.DB)

	hasher := security.NewBcryptHasher()
	//tokenService := service.NewJWTSecret([]byte(gcfg.JWT_SECRET))
	tokenService := service2.NewJwtAuthenticator(gcfg.JWT_SECRET, tokenRepo, userRepo)
	authService := service2.NewAuthService(userRepo, hasher, tokenService, publisher)
	middlewares.Init(tokenService, authService)

	roleRepo := repository.NewRoleRepository(database.DB)
	permissionRepo := repository.NewPermissionRepository(database.DB)
	rbacRepo := repository.NewRBACRepository(database.DB)

	//authService := service2.NewAuthService(userRepo, hasher, tokenService, publisher)
	rbacService := service2.NewRBACService(userRepo, roleRepo, permissionRepo, rbacRepo)
	permissionService := service2.NewPermissionService(permissionRepo)
	roleService := service2.NewRoleService(roleRepo)

	authHandler := grpc.NewAuthHandler(authService)
	rbacHandler := grpc.NewRabcHandler(rbacService)
	permissionHandler := grpc.NewPermissionHandler(permissionService)
	roleHandler := grpc.NewRoleHandler(roleService)

	Server = grpc2.NewServer()
	authpb.RegisterAuthServiceServer(Server, authHandler)
	rbacpb.RegisterRBACServiceServer(Server, rbacHandler)
	permissionpb.RegisterPermissionServiceServer(Server, permissionHandler)
	rolepb.RegisterRoleServiceServer(Server, roleHandler)

}

func GetServer() *grpc2.Server {
	return Server
}
