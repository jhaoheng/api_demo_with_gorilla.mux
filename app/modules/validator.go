package modules

import (
	"regexp"

	"github.com/go-playground/validator"
)

var Validator *validator.Validate

func InitValidate() {
	Validator = validator.New()
	if err := Validator.RegisterValidation("is_allow_password", is_allow_password); err != nil {
		panic(err)
	}
	if err := Validator.RegisterValidation("check_regex", check_regex); err != nil {
		panic(err)
	}
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

func check_regex(field_level validator.FieldLevel) bool {
	var re = regexp.MustCompile(`^[A-Za-z0-9]{1,10}$`)
	return re.MatchString(field_level.Field().String())
}
