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
	w                 http.ResponseWriter
	r                 *http.Request
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
	body := &CreateUserBody{}
	//
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(body)
	if err != nil {
		modules.NewResp(w, r).Set(modules.RespContect{Error: err, Stutus: http.StatusBadRequest})
		return
	}

	//
	api := CreateUser{
		w:                 w,
		r:                 r,
		path:              &CreateUserPath{},
		body:              body,
		model_create_user: models.NewUser(),
	}
	status, err := api.do()
	if err != nil {
		modules.NewResp(w, r).Set(modules.RespContect{
			Error:  err,
			Stutus: status,
		})
		return
	}

	//
	payload := CreateUserResult{}
	modules.NewResp(w, r).Set(modules.RespContect{
		Data:   payload,
		Stutus: http.StatusOK,
	})
}

func (api *CreateUser) do() (int, error) {
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
	err := api.model_create_user.SetAcct(api.body.Account).
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
