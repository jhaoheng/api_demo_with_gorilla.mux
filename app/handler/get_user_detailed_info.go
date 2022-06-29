package handler

import (
	"app/models"
	"app/modules"
	"net/http"

	"github.com/gorilla/context"
)

type GetUserDetailed struct {
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
		access_account: context.Get(r, "account").(string),
		model_get_user: models.NewUser(),
		path:           &GetUserDetailedPath{},
		body:           &GetUserDetailedBody{},
	}
	payload, status, err := api.do(w, r)
	modules.NewResp(w, r).Set(modules.RespContect{
		Data:   payload,
		Error:  err,
		Stutus: status,
	})
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
