package handler

import (
	"net/http"

	"api_demo_with_gorilla.mux/app/models"
	"api_demo_with_gorilla.mux/app/modules"
)

type GetUserDetailed struct {
	access_account string
	model_get_user models.IUser
}

type GetUserDetailedResp struct {
	Account   string `json:"account"`
	Fullname  string `json:"fullname"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewGetUserDetailed(mock_api *GetUserDetailed) func(w http.ResponseWriter, r *http.Request) {
	api := GetUserDetailed{}
	if mock_api == nil {
		api = GetUserDetailed{
			model_get_user: models.NewUser(),
		}
	} else {
		api = *mock_api
	}
	return func(w http.ResponseWriter, r *http.Request) {
		api.access_account = r.Context().Value("account").(string)
		payload, status, err := api.do(w, r)
		modules.NewResp(w, r).Set(modules.RespContect{
			Data:   payload,
			Error:  err,
			Stutus: status,
		})
	}
}

func (api *GetUserDetailed) do(w http.ResponseWriter, r *http.Request) (*GetUserDetailedResp, int, error) {
	result, err := api.model_get_user.SetAcct(api.access_account).Get()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	resp := GetUserDetailedResp{
		Account:   result.Acct,
		Fullname:  result.Fullname,
		CreatedAt: result.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: result.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return &resp, http.StatusOK, nil
}
