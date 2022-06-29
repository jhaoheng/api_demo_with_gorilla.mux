package handler

import (
	"app/modules"
	"net/http"
)

type GetCSRFToken struct{}

func NewGetCSRFToken(mock_api *GetCSRFToken) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		modules.NewResp(w, r).Set(modules.RespContect{
			Stutus: http.StatusOK,
		})
	}
}
