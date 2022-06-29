package handler

import (
	"app/models"
	"app/modules"
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
)

type UpdateUser struct {
	body              *UpdateUserBody
	access_account    string
	model_update_user models.IUser
	model_get_user    models.IUser
}

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

func NewUpdateUser(mock_api *UpdateUser) func(w http.ResponseWriter, r *http.Request) {
	api := UpdateUser{}
	if mock_api == nil {
		api = UpdateUser{
			model_update_user: models.NewUser(),
			model_get_user:    models.NewUser(),
		}
	} else {
		api = *mock_api
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println("=====>", r.Context().Value("account"))
		// fmt.Println("====>", context.Get(r, "account"))
		api.body = &UpdateUserBody{}
		api.access_account = context.Get(r, "account").(string)
		resp, status, err := api.do(w, r)
		modules.NewResp(w, r).Set(modules.RespContect{Data: resp, Error: err, Stutus: status})
	}
}

func (api *UpdateUser) do(w http.ResponseWriter, r *http.Request) (*UpdateUserResp, int, error) {
	api.body = &UpdateUserBody{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(api.body); err != nil {
		return nil, http.StatusBadRequest, err
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
		CreatedAt: result.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: result.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return &payload, http.StatusOK, nil
}
