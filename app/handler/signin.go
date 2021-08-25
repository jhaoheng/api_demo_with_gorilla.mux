package handler

import (
	"app/config"
	"app/models"
	"app/modules"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type SigninObj struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

func Signin(w http.ResponseWriter, r *http.Request) {
	signinObj := SigninObj{}
	//
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&signinObj)
	if err != nil || len(signinObj.Account) == 0 || len(signinObj.Password) == 0 {
		err = errors.New("parameters error")
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		signinFialNotify(err, r.RequestURI, signinObj.Account)
		return
	} else {
		if err := modules.CheckRegex(signinObj.Account, signinObj.Password); err != nil {
			modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
			signinFialNotify(err, r.RequestURI, signinObj.Account)
			return
		}
	}

	//
	user := models.USER{}
	if err := user.GetUserDetail(signinObj.Account); err != nil {
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		signinFialNotify(err, r.RequestURI, signinObj.Account)
		return
	}
	//
	hashPwd := modules.HashPasswrod(signinObj.Password)
	if !strings.EqualFold(hashPwd, user.Pwd) {
		err := errors.New("password error")
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		signinFialNotify(err, r.RequestURI, signinObj.Account)
		return
	}

	//
	respContent := map[string]string{
		"token": modules.NewJWTSrv(config.CFG.JWT_PUBLIC_KEY_PATH, config.CFG.JWT_PRIVATE_KEY_PATH).Encrtpying(signinObj.Account),
	}
	modules.NewResp(w, r).SetSuccess(respContent)
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
