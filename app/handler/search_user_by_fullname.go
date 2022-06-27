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
	user := models.NewUser()
	result, err := user.SetFullname(searchUserByFullnameObj.Fullname).Get()
	if err != nil {
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		return
	}

	data := map[string]string{
		"account":    result.Acct,
		"fullname":   result.Fullname,
		"create_at":  result.CreatedAt.String(),
		"updated_at": result.UpdatedAt.String(),
	}
	modules.NewResp(w, r).SetSuccess(data)
}
