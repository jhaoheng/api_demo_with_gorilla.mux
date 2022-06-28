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
	model_get_user    models.IUser
}

type CreateUserPath struct {
}

type CreateUserBody struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Fullname string `json:"fullname"`
}

type CreateUserResp struct {
	Account   string `json:"account"`
	Fullname  string `json:"fullname"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	api := CreateUser{
		path:              &CreateUserPath{},
		body:              &CreateUserBody{},
		model_create_user: models.NewUser(),
		model_get_user:    models.NewUser(),
	}
	payload, status, err := api.do(w, r)
	//
	modules.NewResp(w, r).Set(modules.RespContect{
		Data:   payload,
		Stutus: status,
		Error:  err,
	})
}

func (api *CreateUser) do(w http.ResponseWriter, r *http.Request) (*CreateUserResp, int, error) {
	//
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&api.body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	//
	if len(api.body.Account) == 0 ||
		len(api.body.Password) == 0 ||
		len(api.body.Fullname) == 0 {
		err := fmt.Errorf("parameters lost")
		return nil, http.StatusUnprocessableEntity, err
	}

	//
	if err := modules.CheckRegex(api.body.Account, api.body.Password, api.body.Fullname); err != nil {
		return nil, http.StatusBadRequest, err
	}

	//
	err = api.model_create_user.SetAcct(api.body.Account).
		SetFullname(api.body.Fullname).
		SetPwd(modules.HashPasswrod(api.body.Password)).Create()
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, http.StatusBadRequest, err
		}
		return nil, http.StatusBadGateway, err
	}

	//
	api.model_get_user.SetAcct(api.body.Account)
	api.model_get_user.SetFullname(api.body.Fullname)
	result, _ := api.model_get_user.Get()

	payload := CreateUserResp{
		Account:   result.Acct,
		Fullname:  result.Fullname,
		CreatedAt: result.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: result.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return &payload, 200, nil
}
