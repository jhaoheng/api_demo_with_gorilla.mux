package handler

import (
	"app/models"
	"app/modules"
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/context"

	"github.com/gorilla/mux"
)

type DeleteUser struct {
	w                 http.ResponseWriter
	r                 *http.Request
	path              *DeleteUserPath
	body              *DeleteUserBody
	access_account    string
	model_del_account models.IUser
}

type DeleteUserPath struct {
	DelAccount string
}
type DeleteUserBody struct{}
type DeleteUserResp struct{}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	api := DeleteUser{
		w: w,
		r: r,
		path: &DeleteUserPath{
			DelAccount: vars["account"],
		},
		body:              &DeleteUserBody{},
		access_account:    context.Get(r, "account").(string),
		model_del_account: models.NewUser(),
	}
	status, err := api.do()
	//
	modules.NewResp(w, r).Set(modules.RespContect{
		Error:  err,
		Stutus: status,
		Data:   &DeleteUserResp{},
	})
}

func (api *DeleteUser) do() (int, error) {

	//
	if err := modules.CheckRegex(api.path.DelAccount); err != nil {
		return http.StatusBadRequest, err
	}
	//

	if strings.EqualFold(api.access_account, api.path.DelAccount) {
		err := errors.New("you can't delte yourself")
		return http.StatusBadRequest, err
	}

	//
	api.model_del_account.SetAcct(api.path.DelAccount)
	if rowsAffected, err := api.model_del_account.Delete(); err != nil {
		return http.StatusBadRequest, err
	} else if rowsAffected == 0 {
		err = errors.New("delete account not exist")
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}
