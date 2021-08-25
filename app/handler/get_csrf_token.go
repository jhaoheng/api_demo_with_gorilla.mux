package handler

import (
	"app/modules"
	"net/http"
)

func GetCSRFToken(w http.ResponseWriter, r *http.Request) {
	modules.NewResp(w, r).SetSuccess("success")
}
