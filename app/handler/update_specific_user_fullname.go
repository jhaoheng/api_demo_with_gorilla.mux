package handler

import (
	"app/models"
	"app/modules"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type UpdateUserFullname struct {
	path              *UpdateUserFullnamePath
	body              *UpdateUserFullnameBody
	model_user_update models.IUser
	model_user_get    models.IUser
}

type UpdateUserFullnamePath struct {
	Account string
}

type UpdateUserFullnameBody struct {
	Fullname string
}

type UpdateUserFullnameResp struct {
	Account   string `json:"account"`
	Fullname  string `json:"fullname"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func UpdateUserFullnameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	body := &UpdateUserFullnameBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(body)
	if err != nil {
		err = errors.New("parameters error")
		modules.NewResp(w, r).Set(modules.RespContect{Error: err, Stutus: http.StatusBadRequest})
		return
	}
	//
	api := UpdateUserFullname{
		path: &UpdateUserFullnamePath{
			Account: vars["account"],
		},
		body:              body,
		model_user_update: models.NewUser(),
		model_user_get:    models.NewUser(),
	}
	resp, status, err := api.do()
	modules.NewResp(w, r).Set(modules.RespContect{Data: resp, Error: err, Stutus: status})
}

func (api *UpdateUserFullname) do() (*UpdateUserFullnameResp, int, error) {
	if err := modules.CheckRegex(api.body.Fullname, api.path.Account); err != nil {
		return nil, http.StatusBadRequest, err
	}

	//
	_, err := api.model_user_update.SetAcct(api.path.Account).Update(models.User{
		Fullname: api.body.Fullname,
	})
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	//
	result, err := api.model_user_get.SetAcct(api.path.Account).Get()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	payload := UpdateUserFullnameResp{
		Account:   result.Acct,
		Fullname:  result.Fullname,
		CreatedAt: result.CreatedAt.String(),
		UpdatedAt: result.UpdatedAt.String(),
	}
	return &payload, http.StatusOK, nil
}
