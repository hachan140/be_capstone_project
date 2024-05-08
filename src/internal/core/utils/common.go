package utils

import (
	"be-capstone-project/src/internal/core/common"
	"regexp"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numberBytes = "1234567890"
const VietnamPhoneNumberPrefix = "+84"

var (
	emailRegex, _ = regexp.Compile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	phoneRegex, _ = regexp.Compile("^0[1-9]\\d{8}$")
)

func ValidatePhoneNumber(phoneNumber, autoPhonePattern string) (bool, bool) {
	isOk, err := regexp.MatchString(`^\+84[1-9]\d{8}$`, phoneNumber)
	if err != nil {
		return false, false
	}
	matched, err1 := regexp.Match(autoPhonePattern, []byte(phoneNumber))
	if err1 != nil {
		return false, false
	}

	return isOk, matched
}

func IsValidPhoneNumber(phoneNumber string) bool {
	isOk, err := regexp.MatchString(`^\+84[1-9]\d{8}$`, phoneNumber)
	if err != nil {
		return false
	}

	return isOk
}

func ToGlobalPhoneNumber(strPhoneNumber string) string {
	if !strings.HasPrefix(strPhoneNumber, common.VietnamPhoneNumberPrefix) && strings.HasPrefix(strPhoneNumber, "0") {
		strPhoneNumber = strings.Replace(strPhoneNumber, "0", common.VietnamPhoneNumberPrefix, 1)
	}

	return strPhoneNumber
}

func ToLocalPhoneNumber(strPhoneNumber string) string {
	if strings.HasPrefix(strPhoneNumber, VietnamPhoneNumberPrefix) && !strings.HasPrefix(strPhoneNumber, "0") {
		strPhoneNumber = strings.Replace(strPhoneNumber, VietnamPhoneNumberPrefix, "0", 1)
	}

	return strPhoneNumber
}

func ValidateEmail(email string) bool {
	//re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(email)
}
