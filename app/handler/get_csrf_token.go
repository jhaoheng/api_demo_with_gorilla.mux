package handler

import (
	"net/http"

	"api_demo_with_gorilla.mux/app/modules"
)

type GetCSRFToken struct{}

func NewGetCSRFToken(mock_api *GetCSRFToken) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		modules.NewResp(w, r).Set(modules.RespContect{
			Stutus: http.StatusOK,
		})
	}
}
