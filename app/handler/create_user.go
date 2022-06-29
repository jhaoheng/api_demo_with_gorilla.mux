package handler

import (
	"app/models"
	"app/modules"
	"encoding/json"
	"net/http"
)

type CreateUser struct {
	body              *CreateUserBody
	model_create_user models.IUser
	model_get_user    models.IUser
}

type CreateUserBody struct {
	Account  string `json:"account" validate:"required,check_regex"`
	Password string `json:"password" validate:"required,is_allow_password"`
	Fullname string `json:"fullname" validate:"required,check_regex"`
}

type CreateUserResp struct {
	Account   string `json:"account"`
	Fullname  string `json:"fullname"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewCreateUser(mock_api *CreateUser) func(w http.ResponseWriter, r *http.Request) {
	api := CreateUser{}
	if mock_api == nil {
		api = CreateUser{
			body:              &CreateUserBody{},
			model_create_user: models.NewUser(),
			model_get_user:    models.NewUser(),
		}
	} else {
		api = *mock_api
	}
	return func(w http.ResponseWriter, r *http.Request) {
		payload, status, err := api.do(w, r)
		//
		modules.NewResp(w, r).Set(modules.RespContect{
			Data:   payload,
			Stutus: status,
			Error:  err,
		})
	}
}

func (api *CreateUser) do(w http.ResponseWriter, r *http.Request) (*CreateUserResp, int, error) {
	//
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&api.body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	//
	if err := modules.Validate(api.body); err != nil {
		return nil, http.StatusUnprocessableEntity, err
	}

	//
	err = api.model_create_user.SetAcct(api.body.Account).
		SetFullname(api.body.Fullname).
		SetPwd(modules.HashPasswrod(api.body.Password)).Create()
	if err != nil {
		return nil, http.StatusBadRequest, err
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
