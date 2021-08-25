package handler

import (
	"app/models"
	"app/modules"
	"net/http"

	"github.com/gorilla/mux"
)

type SearchUserByFullnameObj struct {
	Fullname string
}

func SearchUserByFullname(w http.ResponseWriter, r *http.Request) {
	searchUserByFullnameObj := SearchUserByFullnameObj{}
	//
	vars := mux.Vars(r)
	searchUserByFullnameObj.Fullname = vars["fullname"]

	//
	if err := modules.CheckRegex(searchUserByFullnameObj.Fullname); err != nil {
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		return
	}

	//
	user := models.USER{}
	if err := user.SearchByFullname(searchUserByFullnameObj.Fullname); err != nil {
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		return
	}

	data := map[string]string{
		"account":    user.Acct,
		"fullname":   user.Fullname,
		"create_at":  user.CreatedAt.String(),
		"updated_at": user.UpdatedAt.String(),
	}
	modules.NewResp(w, r).SetSuccess(data)
}
