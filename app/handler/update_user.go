package handler

import (
	"app/models"
	"app/modules"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/context"
)

type UpdateUserObj struct {
	Password string `json:"password"`
	Fullname string `json:"fullname"`
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	obj := UpdateUserObj{}
	//
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&obj)
	if err != nil || (len(obj.Fullname) == 0 && len(obj.Password) == 0) {
		err = errors.New("parameters error")
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		return
	} else {
		if len(obj.Fullname) != 0 {
			if err := modules.CheckRegex(obj.Fullname); err != nil {
				modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
				return
			}
		}
		if len(obj.Password) != 0 {
			if err := modules.CheckRegex(obj.Password); err != nil {
				modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
				return
			}
		}
	}

	//
	account := context.Get(r, "account").(string)
	user := models.USER{
		Acct: account,
	}
	if _, err := user.Update(modules.HashPasswrod(obj.Password), obj.Fullname); err != nil {
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		return
	}
	user.GetUserDetail(account)

	//
	data := map[string]string{
		"account":    user.Acct,
		"fullname":   user.Fullname,
		"create_at":  user.CreatedAt.String(),
		"updated_at": user.UpdatedAt.String(),
	}
	modules.NewResp(w, r).SetSuccess(data)
}
