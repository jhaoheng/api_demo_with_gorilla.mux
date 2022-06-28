package handler

import (
	"app/models"
	"app/modules"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type CreateUser struct {
	path              *CreateUserPath
	body              *CreateUserBody
	model_create_user models.IUser
}

type CreateUserPath struct {
}

type CreateUserBody struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Fullname string `json:"fullname"`
}

type CreateUserResult struct {
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	api := CreateUser{
		path:              &CreateUserPath{},
		body:              &CreateUserBody{},
		model_create_user: models.NewUser(),
	}
	status, err := api.do(w, r)
	//
	modules.NewResp(w, r).Set(modules.RespContect{
		Data:   CreateUserResult{},
		Stutus: status,
		Error:  err,
	})
}

func (api *CreateUser) do(w http.ResponseWriter, r *http.Request) (int, error) {
	//
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&api.body)
	if err != nil {
		return http.StatusBadRequest, err
	}

	//
	if len(api.body.Account) == 0 ||
		len(api.body.Password) == 0 ||
		len(api.body.Fullname) == 0 {
		err := fmt.Errorf("parameters lost")
		return http.StatusUnprocessableEntity, err
	}

	//
	if err := modules.CheckRegex(api.body.Account, api.body.Password, api.body.Fullname); err != nil {
		return http.StatusBadRequest, err
	}

	//
	err = api.model_create_user.SetAcct(api.body.Account).
		SetFullname(api.body.Fullname).
		SetPwd(modules.HashPasswrod(api.body.Password)).Create()
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return http.StatusBadRequest, err
		}
		return http.StatusBadGateway, err
	}
	return 200, nil
}
