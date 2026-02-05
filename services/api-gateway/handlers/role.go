package handlers

import (
	rolepb "abbassmortazavi/go-microservice/pkg/proto/role"
	"abbassmortazavi/go-microservice/pkg/utils"
	"abbassmortazavi/go-microservice/services/api-gateway/grpc_clients"
	"abbassmortazavi/go-microservice/services/api-gateway/requests/role"
	"net/http"
	"strconv"
)

func Create(w http.ResponseWriter, r *http.Request) {
	var req role.CreateRoleReq
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
	res, err := authService.Role.Create(ctx, req.ToProto())
	if err != nil {
		utils.InternalError(w, err)
		return
	}
	utils.Success(w, http.StatusOK, res, "Role Has been inserted Successfully!!")
}
func Delete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetPathParamInt(r, "id")
	if err != nil {
		utils.BadRequest(w, "invalid id", err)
		return
	}
	var req role.DeleteRoleReq
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
	_, err = authService.Role.Delete(ctx, req.ToProto())
	if err != nil {
		utils.InternalError(w, err)
		return
	}
	utils.Success(w, http.StatusOK, "", "Role has been deleted Successfully!")
}
func List(w http.ResponseWriter, r *http.Request) {
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
	req := &rolepb.ListRolesRequest{
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
	res, err := authService.Role.Lists(ctx, req)
	if err != nil {
		utils.InternalError(w, err)
		return
	}

	utils.Success(w, http.StatusOK, res, "Permission has been listed Successfully!")
}

func Get(w http.ResponseWriter, r *http.Request) {
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
	req := rolepb.GetRoleDetailRequest{
		Id: int64(id),
	}
	res, err := authService.Role.Get(ctx, &req)
	if err != nil {
		utils.InternalError(w, err)
		return
	}
	utils.Success(w, http.StatusOK, res, "Role has been retrieved Successfully!")
}
