package models

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	const (
		APP_DB_USERNAME = "ui_test"
		APP_DB_PASSWORD = "ui_test"
		APP_DB_NAME     = "ui_test"
	)
	NewDBConnection("localhost", APP_DB_USERNAME, APP_DB_PASSWORD, APP_DB_NAME)
	m.Run()
}

func Test_User_ListAll(t *testing.T) {
	var total int64 = 0
	var err error
	users := USERS{}
	if err = users.ListBy("0", "asc", 10, &total); err != nil {
		t.Fatal(err)
	}
	fmt.Println("Count =>", total)
	for _, user := range users {
		log.Println(user)
	}
}

func Test_User_SearchByFullname(t *testing.T) {
	user := USER{}
	if err := user.SearchByFullname("maxhu11111"); err != nil {
		t.Fatal(err)
	}
	log.Println(user)
}

func Test_User_GetUserDetail(t *testing.T) {
	user := USER{}
	if err := user.GetUserDetail("account_696327000"); err != nil {
		t.Fatal(err)
	}
	log.Println(user)
}

func Test_User_Create(t *testing.T) {
	tmp := time.Now().Nanosecond()
	user := USER{
		Acct:     fmt.Sprintf("account_%v", tmp),
		Pwd:      "ui_test",
		Fullname: "fullname",
	}
	if err := user.Create(); err != nil {
		t.Fatal(err)
	}
}

func Test_User_Delete(t *testing.T) {
	user := USER{
		Acct: "account_696821000",
	}
	if _, err := user.Delete(); err != nil {
		t.Fatal(err)
	}
}

func Test_User_Update(t *testing.T) {
	user := USER{
		Acct: "account_696327000",
	}
	if rowsAffected, err := user.Update("password1", "fullname"); err != nil {
		t.Fatal(err)
	} else {
		log.Println("rowsAffected =>", rowsAffected)
	}
}
