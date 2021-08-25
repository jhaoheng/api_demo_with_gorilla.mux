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

type DeleteUserObj struct {
	DelAccount string
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	deleteUserObj := DeleteUserObj{}
	//
	vars := mux.Vars(r)
	deleteUserObj.DelAccount = vars["account"]
	if err := modules.CheckRegex(deleteUserObj.DelAccount); err != nil {
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		return
	}

	//
	account := context.Get(r, "account").(string)
	if strings.EqualFold(account, deleteUserObj.DelAccount) {
		err := errors.New("you can't delte yourself")
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		return
	}

	user := models.USER{
		Acct: deleteUserObj.DelAccount,
	}
	if rowsAffected, err := user.Delete(); err != nil {
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		return
	} else if rowsAffected == 0 {
		err = errors.New("delete account not exist")
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		return
	}

	modules.NewResp(w, r).SetSuccess("success")
}
