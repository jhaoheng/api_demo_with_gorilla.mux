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
	user := models.NewUser()
	user.SetAcct(updateObj.Account)
	rowsAffected, err := user.Update(models.User{
		Fullname: updateObj.Fullname,
	})
	if err != nil {
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		return
	} else if rowsAffected == 0 {
		err = errors.New("specific account is not exist")
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		return
	}

	//
	user = models.NewUser()
	result, err := user.SetAcct(updateObj.Account).Get()
	if err != nil {
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		return
	}

	//
	data := map[string]string{
		"account":    result.Acct,
		"fullname":   result.Fullname,
		"create_at":  result.CreatedAt.String(),
		"updated_at": result.UpdatedAt.String(),
	}
	modules.NewResp(w, r).SetSuccess(data)
}
