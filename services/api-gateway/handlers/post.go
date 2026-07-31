package handlers

import (
	"abbassmortazavi/go-microservice/pkg/utils"
	"abbassmortazavi/go-microservice/services/api-gateway/grpc_clients"
	"context"
	"log"
	"net/http"

	"google.golang.org/protobuf/types/known/emptypb"
)

func PostList(w http.ResponseWriter, r *http.Request) {
	log.Println("PostHandler.List")
	postService, err := grpc_clients.NewPostClientService()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := postService.Client.List(context.Background(), &emptypb.Empty{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = utils.WriteJson(w, http.StatusOK, res)
	if err != nil {
		return
	}
}
