package handlers

import (
	permissionpb "abbassmortazavi/go-microservice/pkg/proto/abbassmortazavi/go-microservice/permission"
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
	res, err := authService.Permission.Create(ctx, req.ToProto())
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
	_, err = authService.Permission.Delete(ctx, req.ToProto())
	if err != nil {
		utils.InternalError(w, err)
		return
	}
	utils.Success(w, http.StatusOK, "", "Permission has been deleted Successfully!")
}
func ListPermissions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	ctx := r.Context()
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.ParseInt(query.Get("per_page"), 10, 64)
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	search := query.Get("search")
	sortBy := query.Get("sort_by")
	orderBy := query.Get("order_by")

	if sortBy == "" {
		sortBy = "desc"
	}
	if orderBy == "" {
		orderBy = "id"
	}

	// Create protobuf request
	req := &permissionpb.ListPermissionsRequest{
		Page:    int64(page),
		PerPage: perPage,
		Search:  search,
		SortBy:  sortBy,
		OrderBy: orderBy,
	}

	authService, err := grpc_clients.NewAuthServiceClient()
	if err != nil {
		utils.InternalError(w, err)
		return
	}
	res, err := authService.Permission.Lists(ctx, req)
	if err != nil {
		utils.InternalError(w, err)
		return
	}

	utils.Success(w, http.StatusOK, res, "Permission has been listed Successfully!")
}

func GetPermission(w http.ResponseWriter, r *http.Request) {
	log.Println("GetPermission")
	id, err := utils.GetPathParamInt(r, "id")
	if err != nil {
		utils.BadRequest(w, "invalid id", err)
		return
	}
	authService, err := grpc_clients.NewAuthServiceClient()
	if err != nil {
		utils.InternalError(w, err)
		return
	}
	ctx := r.Context()
	req := permissionpb.GetPermissionDetailRequest{
		Id: int64(id),
	}
	res, err := authService.Permission.Get(ctx, &req)
	if err != nil {
		utils.InternalError(w, err)
		return
	}
	utils.Success(w, http.StatusOK, res, "Permission has been retrieved Successfully!")
}
