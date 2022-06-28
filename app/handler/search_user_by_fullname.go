package handler

import (
	"app/models"
	"app/modules"
	"net/http"

	"github.com/gorilla/mux"
)

type SearchUserByFullname struct {
	path           *SearchUserByFullnamePath
	body           *SearchUserByFullnameBody
	model_get_user models.IUser
}

type SearchUserByFullnamePath struct {
	Fullname string
}
type SearchUserByFullnameBody struct{}
type SearchUserByFullnameResp struct {
	Account   string `json:"account"`
	Fullname  string `json:"fullname"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func SearchUserByFullnameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	api := SearchUserByFullname{
		path: &SearchUserByFullnamePath{
			Fullname: vars["fullname"],
		},
		model_get_user: models.NewUser(),
	}
	resp, status, err := api.do()
	modules.NewResp(w, r).Set(modules.RespContect{Data: resp, Error: err, Stutus: status})
}

func (api *SearchUserByFullname) do() (*SearchUserByFullnameResp, int, error) {

	if err := modules.CheckRegex(api.path.Fullname); err != nil {
		return nil, http.StatusBadRequest, err
	}

	result, err := api.model_get_user.SetFullname(api.path.Fullname).Get()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	//
	payload := SearchUserByFullnameResp{
		Account:   result.Acct,
		Fullname:  result.Fullname,
		CreatedAt: result.CreatedAt.String(),
		UpdatedAt: result.UpdatedAt.String(),
	}
	return &payload, http.StatusOK, nil
}
