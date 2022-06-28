package handler

import (
	"app/models"
	"app/modules"
	"net/http"

	"github.com/gorilla/context"
)

type GetUserDetailed struct {
	w              http.ResponseWriter
	r              *http.Request
	path           *GetUserDetailedPath
	body           *GetUserDetailedBody
	access_account string
	model_get_user models.IUser
}

type GetUserDetailedPath struct{}
type GetUserDetailedBody struct{}
type GetUserDetailedResp struct {
	Account   string `json:"account"`
	Fullname  string `json:"fullname"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func GetUserDetailedHandler(w http.ResponseWriter, r *http.Request) {
	api := GetUserDetailed{
		w:              w,
		r:              r,
		access_account: context.Get(r, "account").(string),
		model_get_user: models.NewUser(),
		path:           &GetUserDetailedPath{},
		body:           &GetUserDetailedBody{},
	}
	payload, status, err := api.do()
	modules.NewResp(w, r).Set(modules.RespContect{
		Data:   payload,
		Error:  err,
		Stutus: status,
	})
}

func (api *GetUserDetailed) do() (*GetUserDetailedResp, int, error) {
	result, err := api.model_get_user.SetAcct(api.access_account).Get()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	resp := GetUserDetailedResp{
		Account:   result.Acct,
		Fullname:  result.Fullname,
		CreatedAt: result.CreatedAt.String(),
		UpdatedAt: result.UpdatedAt.String(),
	}
	return &resp, http.StatusOK, nil
}
