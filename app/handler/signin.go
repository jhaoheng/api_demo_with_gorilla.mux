package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"api_demo_with_gorilla.mux/app/config"
	"api_demo_with_gorilla.mux/app/models"
	"api_demo_with_gorilla.mux/app/modules"
)

type Signin struct {
	body                 *SigninBody
	model_get_user       models.IUser
	jwt_public_key_path  string
	jwt_private_key_path string
}

type SigninBody struct {
	Account  string `json:"account" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type SigninResp struct {
	Token string `json:"token"`
}

func NewSignin(mock_api *Signin) func(w http.ResponseWriter, r *http.Request) {
	api := Signin{}
	if mock_api == nil {
		api = Signin{
			model_get_user:       models.NewUser(),
			jwt_public_key_path:  config.CFG.JWT_PUBLIC_KEY_PATH,
			jwt_private_key_path: config.CFG.JWT_PRIVATE_KEY_PATH,
		}
	} else {
		api = *mock_api
	}
	return func(w http.ResponseWriter, r *http.Request) {
		resp, status, err := api.do(w, r)
		signinFialNotify(err, r.RequestURI, api.body.Account)
		modules.NewResp(w, r).Set(modules.RespContect{Data: resp, Stutus: status, Error: err})
	}
}

func (api *Signin) do(w http.ResponseWriter, r *http.Request) (*SigninResp, int, error) {
	api.body = &SigninBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(api.body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	//
	if err := modules.Validate(api.body); err != nil {
		return nil, http.StatusUnprocessableEntity, err
	}

	//
	hashPwd := modules.HashPasswrod(api.body.Password)
	_, err = api.model_get_user.SetAcct(api.body.Account).SetPwd(hashPwd).Get()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	//
	jwt, err := modules.NewJWTSrv(config.JWTPubKey, config.JWTPriKey)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	token := jwt.Encrtpying(api.body.Account)

	//
	resp := SigninResp{
		Token: token,
	}
	return &resp, http.StatusOK, nil
}

func signinFialNotify(err error, from, account string) {
	if WSConnections != 0 {
		go func() {
			wsmsg := WSErrMessage{
				Account: account,
				Err: func() string {
					if err != nil {
						return err.Error()
					}
					return ""
				}(),
				From: from,
				At:   time.Now(),
			}
			WSChannel <- wsmsg
		}()
	}
}
