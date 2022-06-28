package modules

import (
	"regexp"

	"github.com/go-playground/validator"
)

var Validator *validator.Validate

func InitValidate() {
	Validator = validator.New()
	Validator.RegisterValidation("is_allow_password", is_allow_password)
}

func Validate(obj interface{}) error {
	return Validator.Struct(obj)
}

//
func is_allow_password(field_level validator.FieldLevel) bool {
	pattern := "^([A-Za-z0-9]){3,32}$"
	ok, _ := regexp.MatchString(pattern, field_level.Field().String())
	return ok
}
