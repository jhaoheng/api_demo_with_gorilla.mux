package middlewares

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func ShowRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/health":
		default:
			logrus.Info(r.RequestURI)
		}
		next.ServeHTTP(w, r)
	})
}
