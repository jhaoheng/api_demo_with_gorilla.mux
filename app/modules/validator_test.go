package modules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Validator(t *testing.T) {
	type Data struct {
		input   string
		err_str string
	}
	datas := []Data{
		0: {input: "123456789", err_str: ""},
		1: {input: "***123", err_str: "Key: 'User.Password' Error:Field validation for 'Password' failed on the 'is_allow_password' tag"},
		2: {input: "", err_str: "Key: 'User.Password' Error:Field validation for 'Password' failed on the 'required' tag"},
	}

	InitValidate()
	type User struct {
		Password string `validate:"required,is_allow_password"`
	}

	for _, data := range datas {
		user := User{
			Password: data.input,
		}
		err := Validate(user)
		err_str := ""
		if err != nil {
			err_str = err.Error()
		}
		if !assert.Equal(t, data.err_str, err_str) {
			t.Fatal()
		}
	}
}
