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
	user := models.NewUser()
	user.SetAcct(account)
	_, err = user.Update(models.User{
		Pwd:      modules.HashPasswrod(obj.Password),
		Fullname: obj.Fullname,
	})
	if err != nil {
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		return
	}
	//
	user = models.NewUser()
	result, err := user.SetAcct(account).Get()
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
