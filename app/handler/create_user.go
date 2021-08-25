package handler

import (
	"app/models"
	"app/modules"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type CreateUserObj struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Fullname string `json:"fullname"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	createUserObj := CreateUserObj{}
	//
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&createUserObj)
	if err != nil ||
		len(createUserObj.Account) == 0 ||
		len(createUserObj.Password) == 0 ||
		len(createUserObj.Fullname) == 0 {
		err = errors.New("parameters error")
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		return
	} else {
		if err := modules.CheckRegex(createUserObj.Account, createUserObj.Password, createUserObj.Fullname); err != nil {
			modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
			return
		}
	}

	//
	user := models.USER{
		Acct:     createUserObj.Account,
		Pwd:      modules.HashPasswrod(createUserObj.Password),
		Fullname: createUserObj.Fullname,
	}

	if err := user.Create(); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		} else {
			modules.NewResp(w, r).SetError(err, http.StatusBadGateway)
		}
		return
	}

	modules.NewResp(w, r).SetSuccess("success")
}
