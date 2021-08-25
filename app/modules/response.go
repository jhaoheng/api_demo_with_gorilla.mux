package modules

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/csrf"
)

type RESP struct {
	writer  http.ResponseWriter
	request *http.Request
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

type RESPCONTENT struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

func NewResp(w http.ResponseWriter, r *http.Request) *RESP {
	return &RESP{
		writer:  w,
		request: r,
	}
}

func (r *RESP) SetSuccess(data interface{}) {
	w := r.writer
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-CSRF-Token", csrf.Token(r.request))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	respContent := RESPCONTENT{
		Data: data,
	}
	b, _ := json.Marshal(respContent)
	w.Write(b)
}

func (r *RESP) SetError(err error, statusCode int) {
	w := r.writer
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	respContent := RESPCONTENT{
		Error: err.Error(),
	}
	b, _ := json.Marshal(respContent)
	w.Write(b)
}
