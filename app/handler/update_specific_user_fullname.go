package handler

import (
	"app/models"
	"app/modules"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type UpdateSpecificUserFullnameObj struct {
	Account  string
	Fullname string
}

func UpdateSpecificUserFullname(w http.ResponseWriter, r *http.Request) {
	updateObj := UpdateSpecificUserFullnameObj{}
	//
	vars := mux.Vars(r)
	updateObj.Account = vars["account"]

	//
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updateObj); err != nil || len(updateObj.Fullname) == 0 {
		err = errors.New("parameters error")
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		return
	} else {
		if err := modules.CheckRegex(updateObj.Fullname, updateObj.Account); err != nil {
			modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
			return
		}
	}

	//
	user := models.USER{
		Acct: updateObj.Account,
	}
	if rowsAffected, err := user.Update("", updateObj.Fullname); err != nil {
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		return
	} else if rowsAffected == 0 {
		err = errors.New("specific account is not exist")
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		return
	}
	user.GetUserDetail(updateObj.Account)

	data := map[string]string{
		"account":    user.Acct,
		"fullname":   user.Fullname,
		"create_at":  user.CreatedAt.String(),
		"updated_at": user.UpdatedAt.String(),
	}
	modules.NewResp(w, r).SetSuccess(data)
}
