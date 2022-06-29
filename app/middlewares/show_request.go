package middlewares

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func ShowRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/health" {
			return
		}
		logrus.Info(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
