package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type HomeApi struct {
}

func NewHomeApi() *HomeApi {
	return &HomeApi{}
}

func (ha *HomeApi) Register(router *mux.Router) {
	router.HandleFunc("/home", ha.show).Methods("GET")
}

func (ha *HomeApi) show(w http.ResponseWriter, r *http.Request) {

}
