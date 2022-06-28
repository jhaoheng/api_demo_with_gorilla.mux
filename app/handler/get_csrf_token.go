package handler

import (
	"app/modules"
	"net/http"
)

func GetCSRFTokenHandler(w http.ResponseWriter, r *http.Request) {
	modules.NewResp(w, r).Set(modules.RespContect{
		Stutus: http.StatusOK,
	})
}
