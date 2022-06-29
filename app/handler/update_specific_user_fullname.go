package handler

import (
	"app/models"
	"app/modules"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type UpdateUserFullname struct {
	path              *UpdateUserFullnamePath
	body              *UpdateUserFullnameBody
	model_update_user models.IUser
	model_get_user    models.IUser
}

type UpdateUserFullnamePath struct {
	Account string
}

type UpdateUserFullnameBody struct {
	Fullname string `validate:"required"`
}

type UpdateUserFullnameResp struct {
	Account   string `json:"account"`
	Fullname  string `json:"fullname"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewUpdateUserFullname(mock_api *UpdateUserFullname) func(w http.ResponseWriter, r *http.Request) {
	api := UpdateUserFullname{}
	if mock_api == nil {
		api = UpdateUserFullname{
			model_update_user: models.NewUser(),
			model_get_user:    models.NewUser(),
		}
	} else {
		api = *mock_api
	}
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		api.path = &UpdateUserFullnamePath{
			Account: vars["account"],
		}
		api.body = &UpdateUserFullnameBody{}
		resp, status, err := api.do(w, r)
		modules.NewResp(w, r).Set(modules.RespContect{Data: resp, Error: err, Stutus: status})
	}
}

func (api *UpdateUserFullname) do(w http.ResponseWriter, r *http.Request) (*UpdateUserFullnameResp, int, error) {
	//
	api.body = &UpdateUserFullnameBody{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(api.body); err != nil {
		return nil, http.StatusBadRequest, err
	}

	//
	if err := modules.Validate(api.body); err != nil {
		return nil, http.StatusUnprocessableEntity, err
	}

	//
	_, err := api.model_update_user.SetAcct(api.path.Account).Update(models.User{
		Fullname: api.body.Fullname,
	})
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	//
	result, err := api.model_get_user.SetAcct(api.path.Account).Get()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	payload := UpdateUserFullnameResp{
		Account:   result.Acct,
		Fullname:  result.Fullname,
		CreatedAt: result.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: result.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return &payload, http.StatusOK, nil
}
