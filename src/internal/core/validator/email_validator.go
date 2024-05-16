package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func EmailValidatorFunc(fl validator.FieldLevel) bool {
	status, _ := IsValidEmail(fl.Field().String())
	return status
}

func IsValidEmail(email string) (bool, *string) {
	emailRegex, _ := regexp.Compile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	match := emailRegex.MatchString(email)
	if match {
		return true, &email
	}

	return false, nil
}
