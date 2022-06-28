package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_User_Create(t *testing.T) {
	var model IUser = &User{}
	model.SetAcct("account_3")
	model.SetPwd("test")
	model.SetFullname("fullname")
	err := model.Create()
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}
}

func Test_User_Get(t *testing.T) {
	var model IUser = &User{}
	model.SetFullname("fullname")
	result, err := model.Get()
	if assert.NoError(t, err) {
		fmt.Printf("%#v\n", result)
	}
}

func Test_User_GetAllCount(t *testing.T) {
	var model IUser = &User{}
	result, err := model.GetAllCount()
	if assert.NoError(t, err) {
		fmt.Printf("%#v\n", result)
	}
}

func Test_User_Update(t *testing.T) {
	var model IUser = &User{}
	model.SetAcct("account")
	_, err := model.Update(User{
		Fullname: "my_updated",
	})
	assert.NoError(t, err)
}

func Test_User_Delete(t *testing.T) {
	var model IUser = &User{}
	model.SetAcct("account")
	_, err := model.Delete()
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}
}

func Test_User_AndOr(t *testing.T) {
	model := NewUser()
	model.SetAcct("account_1")
	model = model.Or(User{Acct: "account_2"}, User{Acct: "account_3"})

	results, err := model.GetAll()
	if assert.NoError(t, err) {
		for _, result := range results {
			fmt.Printf("%#v\n", result.Acct)
		}
	}
}

func Test_User_Mock_Create(t *testing.T) {
	mock_user := NewMockUser()
	mock_user.On("SetAcct", "")
	mock_user.On("SetPwd", "")
	mock_user.On("SetFullname", "")
	mock_user.On("Create").Return(nil)
}
