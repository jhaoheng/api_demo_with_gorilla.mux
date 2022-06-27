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
	user := models.NewUser()
	result, err := user.SetAcct(getUserDetailedInfoObj.Account).Get()
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
