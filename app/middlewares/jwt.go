package middlewares

import (
	"errors"
	"net/http"

	"api_demo_with_gorilla.mux/app/config"
	"api_demo_with_gorilla.mux/app/models"
	"api_demo_with_gorilla.mux/app/modules"

	"github.com/gorilla/context"
	"github.com/sirupsen/logrus"
)

//
func JWTValidate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			err := errors.New("forbidden, authorization fail")
			modules.NewResp(w, r).Set(modules.RespContect{Error: err, Stutus: http.StatusForbidden})
			return
		}

		//
		jwtsrv, err := modules.NewJWTSrv(config.CFG.JWT_PUBLIC_KEY_PATH, config.CFG.JWT_PRIVATE_KEY_PATH)
		if err != nil {
			modules.NewResp(w, r).Set(modules.RespContect{Error: err, Stutus: http.StatusInternalServerError})
			return
		}
		account, ok := jwtsrv.Validating(tokenString)
		if ok {
			logrus.Infof("Authenticated user =>%s\n", account)
			context.Set(r, "account", account) // use => account := context.Get(r, "account").(string)
		} else {
			err := errors.New("forbidden, authorization fail")
			modules.NewResp(w, r).Set(modules.RespContect{Error: err, Stutus: http.StatusForbidden})
			return
		}
		//
		user := models.NewUser()
		user.SetAcct(account)
		_, err = user.Get()
		if err != nil {
			err = errors.New("forbidden, authorization fail")
			modules.NewResp(w, r).Set(modules.RespContect{Error: err, Stutus: http.StatusForbidden})
			return
		}
		//

		next.ServeHTTP(w, r)
	})
}
