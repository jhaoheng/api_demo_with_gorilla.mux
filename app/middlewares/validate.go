package middlewares

import (
	"app/config"
	"app/models"
	"app/modules"
	"errors"
	"net/http"

	"github.com/gorilla/context"
	"github.com/sirupsen/logrus"
)

//
func JWTValidate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			err := errors.New("forbidden, authorization fail")
			modules.NewResp(w, r).SetError(err, http.StatusForbidden)
			return
		}

		//
		jwtsrv := modules.NewJWTSrv(config.CFG.JWT_PUBLIC_KEY_PATH, config.CFG.JWT_PRIVATE_KEY_PATH)
		account, ok := jwtsrv.Validating(tokenString)
		if ok {
			logrus.Infof("Authenticated user =>%s\n", account)
			context.Set(r, "account", account) // use => account := context.Get(r, "account").(string)
		} else {
			err := errors.New("forbidden, authorization fail")
			modules.NewResp(w, r).SetError(err, http.StatusForbidden)
			return
		}
		//
		user := models.NewUser()
		user.SetAcct(account)
		_, err := user.Get()
		if err != nil {
			err = errors.New("forbidden, authorization fail")
			modules.NewResp(w, r).SetError(err, http.StatusForbidden)
			return
		}
		//

		next.ServeHTTP(w, r)
	})
}
