package handler

import (
	"app/config"
	"app/models"
	"app/modules"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Signin struct {
	path           *SigninPath
	body           *SigninBody
	model_get_user models.IUser
}

type SigninPath struct{}
type SigninBody struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}
type SigninResp struct {
	Token string `json:"token"`
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	//
	body := &SigninBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(body)
	if err != nil {
		modules.NewResp(w, r).Set(modules.RespContect{Error: err, Stutus: http.StatusBadRequest})
		signinFialNotify(err, r.RequestURI, body.Account)
		return
	}

	api := Signin{
		path:           &SigninPath{},
		body:           body,
		model_get_user: models.NewUser(),
	}
	resp, status, err := api.do()
	if err != nil {
		signinFialNotify(err, r.RequestURI, api.body.Account)
	}
	modules.NewResp(w, r).Set(modules.RespContect{Data: resp, Stutus: status, Error: err})
}

func (api *Signin) do() (*SigninResp, int, error) {
	if len(api.body.Account) == 0 || len(api.body.Password) == 0 {
		err := fmt.Errorf("parameters lost")
		return nil, http.StatusUnprocessableEntity, err
	}

	//
	if err := modules.CheckRegex(api.body.Account, api.body.Password); err != nil {
		return nil, http.StatusBadRequest, err
	}

	//
	result, err := api.model_get_user.SetAcct(api.body.Account).Get()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	//
	hashPwd := modules.HashPasswrod(api.body.Password)
	if !strings.EqualFold(hashPwd, result.Pwd) {
		err := errors.New("password error")
		return nil, http.StatusBadRequest, err
	}
	//
	resp := SigninResp{
		Token: modules.NewJWTSrv(config.CFG.JWT_PUBLIC_KEY_PATH, config.CFG.JWT_PRIVATE_KEY_PATH).Encrtpying(api.body.Account),
	}
	return &resp, http.StatusOK, nil
}

func signinFialNotify(err error, from, account string) {
	if WSConnections != 0 {
		go func() {
			wsmsg := WSErrMessage{
				Account: account,
				Err:     err.Error(),
				From:    from,
				At:      time.Now(),
			}
			WSChannel <- wsmsg
		}()
	}
}
