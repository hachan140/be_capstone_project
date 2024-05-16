package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

// PhoneValidatorFunc phone validator
func PhoneValidatorFunc(fl validator.FieldLevel) bool {
	status, _ := IsValidPhoneNumber(fl.Field().String())
	return status
}

// IsValidPhoneNumber validate phone number
// - remove all character not off digit
// - allow 8 to 13 digit only
// Input: phone string
// Expect: True if valid, False on else
func IsValidPhoneNumber(phone string) (bool, *string) {
	matchStartZero, _ := regexp.MatchString(`^\+84[1-9]\d{8}$`, phone)
	if !matchStartZero {
		return false, nil
	}

	return true, &phone
}
