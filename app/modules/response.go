package modules

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/csrf"
)

type IResp interface{}

type Resp struct {
	writer  http.ResponseWriter
	request *http.Request
}

type RespContect struct {
	Data   interface{} `json:"data"`
	Error  error       `json:"error"`
	Stutus int         `json:"-"`
}

func NewResp(w http.ResponseWriter, r *http.Request) *Resp {
	return &Resp{
		writer:  w,
		request: r,
	}
}

func (r *Resp) Set(resp RespContect) {
	w := r.writer
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-CSRF-Token", csrf.Token(r.request))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Stutus)
	b, _ := json.Marshal(resp)
	w.Write(b)
}
