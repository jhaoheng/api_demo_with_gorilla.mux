package handler

import (
	"app/models"
	"app/modules"
	"net/http"

	"github.com/gorilla/context"
)

type GetUserDetailedInfoObj struct {
	Account string
}

func GetUserDetailedInfo(w http.ResponseWriter, r *http.Request) {
	getUserDetailedInfoObj := GetUserDetailedInfoObj{}
	getUserDetailedInfoObj.Account = context.Get(r, "account").(string)

	//
	user := models.USER{}
	if err := user.GetUserDetail(getUserDetailedInfoObj.Account); err != nil {
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
