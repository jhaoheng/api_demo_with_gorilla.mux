package handler

import (
	"app/models"
	"app/modules"
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
)

type UpdateUser struct {
	path              *UpdateUserPath
	body              *UpdateUserBody
	access_account    string
	model_update_user models.IUser
	model_get_user    models.IUser
}

type UpdateUserPath struct{}
type UpdateUserBody struct {
	Password string `json:"password"`
	Fullname string `json:"fullname"`
}
type UpdateUserResp struct {
	Account   string `json:"account"`
	Fullname  string `json:"fullname"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	body := UpdateUserBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		modules.NewResp(w, r).Set(modules.RespContect{Error: err, Stutus: http.StatusBadRequest})
		return
	}
	api := UpdateUser{
		path:           &UpdateUserPath{},
		body:           &body,
		access_account: context.Get(r, "account").(string),
	}
	resp, status, err := api.do()
	modules.NewResp(w, r).Set(modules.RespContect{Data: resp, Error: err, Stutus: status})
}

func (api *UpdateUser) do() (*UpdateUserResp, int, error) {

	if len(api.body.Fullname) != 0 {
		if err := modules.CheckRegex(api.body.Fullname); err != nil {
			return nil, http.StatusBadRequest, err
		}
	}
	if len(api.body.Password) != 0 {
		if err := modules.CheckRegex(api.body.Password); err != nil {
			return nil, http.StatusBadRequest, err
		}
	}

	//
	api.model_update_user.SetAcct(api.access_account)
	_, err := api.model_update_user.Update(models.User{
		Pwd:      modules.HashPasswrod(api.body.Password),
		Fullname: api.body.Fullname,
	})
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	//
	result, err := api.model_get_user.SetAcct(api.access_account).Get()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	//
	payload := UpdateUserResp{
		Account:   result.Acct,
		Fullname:  result.Fullname,
		CreatedAt: result.CreatedAt.String(),
		UpdatedAt: result.UpdatedAt.String(),
	}
	return &payload, http.StatusOK, nil
}
