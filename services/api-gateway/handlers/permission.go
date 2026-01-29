package handlers

import (
	"abbassmortazavi/go-microservice/pkg/utils"
	"abbassmortazavi/go-microservice/services/api-gateway/grpc_clients"
	"abbassmortazavi/go-microservice/services/api-gateway/requests/permission"
	"log"
	"net/http"
)

func CreatePermission(w http.ResponseWriter, r *http.Request) {
	log.Println("create permission")
	var req permission.CreatePermissionReq
	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.BadRequest(w, "invalid request", err)
		return
	}
	defer r.Body.Close()
	ctx := r.Context()
	authService, err := grpc_clients.NewAuthServiceClient()
	if err != nil {
		utils.InternalError(w, err)
		return
	}
	res, err := authService.Rbac.CreatePermission(ctx, req.ToProto())
	if err != nil {
		utils.InternalError(w, err)
		return
	}

	utils.Created(w, res)
}
func DeletePermission(w http.ResponseWriter, r *http.Request) {
	log.Println("delete permission")
}
