package handlers

import (
	"abbassmortazavi/go-microservice/pkg/utils"
	"abbassmortazavi/go-microservice/services/api-gateway/grpc_clients"
	"abbassmortazavi/go-microservice/services/api-gateway/requests/rbac"
	"log"
	"net/http"
)

func AssignPermissionToRole(w http.ResponseWriter, r *http.Request) {
	var req rbac.CreateAssignPermissionToRoleReq
	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.BadRequest(w, "Validation request!!", err)
		return
	}
	defer r.Body.Close()
	ctx := r.Context()
	authService, err := grpc_clients.NewAuthServiceClient()
	if err != nil {
		utils.InternalError(w, err)
		return
	}
	res, err := authService.Rbac.AssignPermissionToRole(ctx, req.ToProto())
	if err != nil {
		utils.InternalError(w, err)
		return
	}
	utils.Success(w, http.StatusOK, res, "Role Has been inserted Successfully!!")
}

func AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	var req rbac.CreateAssignRoleToUserReq
	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.BadRequest(w, "Validation request!!", err)
		return
	}
	defer r.Body.Close()
	ctx := r.Context()
	authService, err := grpc_clients.NewAuthServiceClient()
	if err != nil {
		utils.InternalError(w, err)
		return
	}
	res, err := authService.Rbac.AssignRoleToUser(ctx, req.ToProto())
	if err != nil {
		utils.InternalError(w, err)
		return
	}
	utils.Success(w, http.StatusOK, res, "Role Has been Assigned To User Successfully!!")
}
func CheckUserHasRole(w http.ResponseWriter, r *http.Request) {
	var req rbac.CheckUserHasRoleReq
	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.BadRequest(w, "Validation request!!", err)
		return
	}
	defer r.Body.Close()
	ctx := r.Context()
	authService, err := grpc_clients.NewAuthServiceClient()
	if err != nil {
		utils.InternalError(w, err)
		return
	}
	res, err := authService.Rbac.CheckUserRole(ctx, req.ToProto())
	if err != nil {
		utils.InternalError(w, err)
		return
	}
	log.Println("route", res)
	utils.Success(w, http.StatusOK, res, "User Has Permission Has been Checked Successfully!!")
}
