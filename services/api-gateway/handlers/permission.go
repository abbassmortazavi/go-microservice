package handlers

import (
	"abbassmortazavi/go-microservice/pkg/utils"
	"abbassmortazavi/go-microservice/services/api-gateway/grpc_clients"
	"abbassmortazavi/go-microservice/services/api-gateway/requests/permission"
	"log"
	"net/http"
	"strconv"
)

func CreatePermission(w http.ResponseWriter, r *http.Request) {
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
	id, err := utils.GetPathParamInt(r, "id")
	if err != nil {
		utils.BadRequest(w, "invalid id", err)
		return
	}
	var req permission.DeletePermissionReq
	ctx := r.Context()
	authService, err := grpc_clients.NewAuthServiceClient()
	if err != nil {
		utils.InternalError(w, err)
		return
	}
	idConvert, err := strconv.ParseInt(strconv.Itoa(id), 10, 64)
	if err != nil {
		utils.BadRequest(w, "invalid request", err)
		return
	}
	req.Id = idConvert
	_, err = authService.Rbac.DeletePermission(ctx, req.ToProto())
	if err != nil {
		log.Println("111111111111111111111111111111111111")
		utils.InternalError(w, err)
		return
	}
	utils.Success(w, http.StatusOK, "", "Permission has been deleted Successfully!")
}
