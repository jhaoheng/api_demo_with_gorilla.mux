package handler

import (
	"app/models"
	"app/modules"
	"net/http"

	"github.com/gorilla/mux"
)

type SearchUserByFullname struct {
	path           *SearchUserByFullnamePath
	model_get_user models.IUser
}

type SearchUserByFullnamePath struct {
	Fullname string
}
type SearchUserByFullnameResp struct {
	Account   string `json:"account"`
	Fullname  string `json:"fullname"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewSearchUserByFullname(mock_api *SearchUserByFullname) func(w http.ResponseWriter, r *http.Request) {
	api := SearchUserByFullname{}
	if mock_api != nil {
		api = *mock_api
	} else {
		api = SearchUserByFullname{
			model_get_user: models.NewUser(),
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		api.path = &SearchUserByFullnamePath{
			Fullname: vars["fullname"],
		}
		resp, status, err := api.do(w, r)
		modules.NewResp(w, r).Set(modules.RespContect{Data: resp, Error: err, Stutus: status})
	}
}

func (api *SearchUserByFullname) do(w http.ResponseWriter, r *http.Request) (*SearchUserByFullnameResp, int, error) {
	result, err := api.model_get_user.SetFullname(api.path.Fullname).Get()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	//
	payload := SearchUserByFullnameResp{
		Account:   result.Acct,
		Fullname:  result.Fullname,
		CreatedAt: result.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: result.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return &payload, http.StatusOK, nil
}
