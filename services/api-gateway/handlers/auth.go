package handlers

import (
	"abbassmortazavi/go-microservice/pkg/utils"
	"abbassmortazavi/go-microservice/services/api-gateway/grpc_clients"
	"abbassmortazavi/go-microservice/services/api-gateway/requests/auth"
	"abbassmortazavi/go-microservice/services/auth-service/pkg/middlewares"
	"encoding/json"
	"log"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var req auth.RegisterReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	ctx := r.Context()

	authService, err := grpc_clients.NewAuthServiceClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := authService.Client.Register(ctx, req.ToProto())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = utils.WriteJson(w, http.StatusOK, res)
	if err != nil {
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	ctx := r.Context()
	authService, err := grpc_clients.NewAuthServiceClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := authService.Client.Login(ctx, req.ToProto())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = utils.WriteJson(w, http.StatusOK, res)
	if err != nil {
		return
	}
	log.Println("login success")
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	//TODO::
	user, err := middlewares.User(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = utils.WriteJson(w, http.StatusOK, user)
	if err != nil {
		return
	}

	//ctx := r.Context()

	// Convert string ID to int64
	/*userID, err := strconv.ParseInt(user.ID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	req := &authpb.GetUserRequest{
		Id: userID,
	}
	authService, err := grpc_clients.NewAuthServiceClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := authService.Client.GetUser(ctx, req)

	if err != nil {
		utils.InternalError(w, err)
		return
	}
	err = utils.WriteJson(w, http.StatusOK, res)
	if err != nil {
		return
	}*/

}
