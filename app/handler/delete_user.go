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
	path              *DeleteUserPath
	access_account    string
	model_del_account models.IUser
}

type DeleteUserPath struct {
	DelAccount string `validate:"required"`
}
type DeleteUserResp struct{}

func NewDeleteUser(mock_api *DeleteUser) func(w http.ResponseWriter, r *http.Request) {
	api := DeleteUser{}
	if mock_api == nil {
		api = DeleteUser{
			model_del_account: models.NewUser(),
		}
	} else {
		api = *mock_api
	}
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		api.path = &DeleteUserPath{
			DelAccount: vars["account"],
		}
		api.access_account = context.Get(r, "account").(string)
		payload, status, err := api.do(w, r)
		//
		modules.NewResp(w, r).Set(modules.RespContect{
			Error:  err,
			Stutus: status,
			Data:   payload,
		})
	}
}

func (api *DeleteUser) do(w http.ResponseWriter, r *http.Request) (*DeleteUserResp, int, error) {
	//
	if err := modules.Validate(api.path); err != nil {
		return nil, http.StatusUnprocessableEntity, err
	}

	//
	if strings.EqualFold(api.access_account, api.path.DelAccount) {
		err := errors.New("you can't delte yourself")
		return nil, http.StatusBadRequest, err
	}

	//
	api.model_del_account.SetAcct(api.path.DelAccount)
	if rowsAffected, err := api.model_del_account.Delete(); err != nil {
		return nil, http.StatusBadRequest, err
	} else if rowsAffected == 0 {
		err = errors.New("delete account not exist")
		return nil, http.StatusBadRequest, err
	}
	return &DeleteUserResp{}, http.StatusOK, nil
}
