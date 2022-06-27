package handler

import (
	"app/models"
	"app/modules"
	"errors"
	"net/http"
	"regexp"
	"strings"
)

type ListAllUsersObj struct {
	Paging  string
	Sorting string
}

func ListAllUsers(w http.ResponseWriter, r *http.Request) {
	listAllUsersObj := ListAllUsersObj{}
	//
	var err error
	listAllUsersObj.Paging = r.URL.Query().Get("paging")
	listAllUsersObj.Sorting = r.URL.Query().Get("sorting")
	if err = listAllUsersObj.check_paging(); err != nil {
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		return
	}
	if err = listAllUsersObj.check_sorting(); err != nil {
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		return
	}

	//
	user := models.NewUser()
	total, err := user.GetAllCount()
	if err != nil {
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		return
	}
	results, err := user.ListBy(listAllUsersObj.Paging, listAllUsersObj.Sorting, 10)
	if err != nil {
		modules.NewResp(w, r).SetError(err, http.StatusBadRequest)
		return
	}

	//
	datas := make([]map[string]string, 0)
	for _, user := range results {
		data := map[string]string{
			"account":    user.Acct,
			"fullname":   user.Fullname,
			"create_at":  user.CreatedAt.String(),
			"updated_at": user.UpdatedAt.String(),
		}
		datas = append(datas, data)
	}
	respContent := map[string]interface{}{
		"total": total,
		"users": datas,
	}

	modules.NewResp(w, r).SetSuccess(respContent)
}

func (l *ListAllUsersObj) check_paging() error {
	if len(l.Paging) == 0 || strings.EqualFold(l.Paging, "0") {
		l.Paging = "1"
		return nil
	}
	var re = regexp.MustCompile(`[0-9]$`)
	if !re.MatchString(l.Paging) {
		return errors.New("paging must be number")
	}
	return nil
}

func (l *ListAllUsersObj) check_sorting() error {
	if len(l.Sorting) == 0 {
		l.Sorting = "asc"
	}

	var re = regexp.MustCompile(`^(asc|desc)$`)
	if !re.MatchString(l.Sorting) {
		return errors.New(`sorting must be 'asc' or 'desc'`)
	}
	return nil
}
