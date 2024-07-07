package common

import (
	"be-capstone-project/src/internal/core/exception"
	"net/http"
	"strings"
)

type ResponseLanguage string

var (
	VietnameseResponse ResponseLanguage = Vietnamese
	EnglishResponse    ResponseLanguage = English
)

func StringToLanguage(c string) ResponseLanguage {
	c = strings.ToLower(c)
	switch c {
	case "vi", "vietnam":
		return VietnameseResponse
	case "en", "english":
		return EnglishResponse
	default:
		return VietnameseResponse
	}
}

// ErrorResponse error response struct
type ErrorResponse struct {
	HTTPCode       int
	ServiceCode    int
	Message        string
	CustomMsgParam *[]interface{}
}

// ErrorCodeMessage error response struct

type ErrorCodeMessage struct {
	HTTPCode    int
	ServiceCode int
	Message     string
}

func MakeCustomErrorResponse(httpCode int, serviceCode int, msgParams ...interface{}) *ErrorResponse {
	return &ErrorResponse{
		HTTPCode:       httpCode,
		ServiceCode:    serviceCode,
		CustomMsgParam: &[]interface{}{msgParams},
	}
}

func MakeCustomErrorMsgResponse(httpCode int, serviceCode int, msg string) *ErrorResponse {
	return &ErrorResponse{
		HTTPCode:    httpCode,
		ServiceCode: serviceCode,
		Message:     msg,
	}
}

// GetErrorResponse get error response from code
func GetErrorResponse(code int, language string) ErrorResponse {
	var lang ResponseLanguage
	switch language {
	case Vietnamese:
		lang = VietnameseResponse
	case English:
		lang = EnglishResponse
	default:
		lang = VietnameseResponse
	}

	if val, ok := errorResponseMap[code]; ok {
		return val[lang]
	}

	// default response
	return ErrorResponse{
		HTTPCode:    http.StatusInternalServerError,
		ServiceCode: code,
		Message:     http.StatusText(http.StatusInternalServerError),
	}
}

func GetError(err error) int {
	code := http.StatusInternalServerError
	if e, ok := err.(exception.Exception); ok {
		code = int(e.Code())
	}
	return code
}

// Error code 400XXX
const (
	ErrCodeInvalidRequest = 400001

	//authentication
	ErrCodeUserNotFound              = 400002
	ErrCodeInvalidName               = 400003
	ErrCodeInvalidEmail              = 400004
	ErrCodeInvalidPhone              = 400005
	ErrCodeEmailHasAlreadyExisted    = 400006
	ErrCodeUsernameHasAlreadyExisted = 400007
	ErrCodeEmailExisted              = 400009
	ErrCodeInvalidUsername           = 400010
	ErrCodeInvalidUser               = 400011
	ErrCodeUserSocialDoesNotExist    = 400012
	ErrCodeRefreshTokenNotFound      = 400013

	//organization
	ErrCodeOrganizationExisted                  = 400014
	ErrCodeInvalidOrganizationName              = 400015
	ErrCodeOrganizationNotExist                 = 400016
	ErrCodeCannotAccessToOrganization           = 400017
	ErrCodeUserAlreadyInOtherOrganization       = 400018
	ErrCodeUserAlreadyRequestCreateOrganization = 400019

	//category
	ErrCodeCategoryNotFound          = 400019
	ErrCodeCategoryExisted           = 400020
	ErrCodeInvalidCategoryName       = 400021
	ErrCodeInvalidStatus             = 400022
	ErrCodeUserDoesNotHavePermission = 400023

	//document
	ErrCodeDocumentNotFound = 400100

	//user
	ErrCodeUserNotInOrganization = 400200
)

// Error code 401XXX
const (
	ErrCodeInvalidPassword = 401001
)

// Error code 500XXX
const (
	ErrCodeInternalError = 500001
)

var errorResponseMap = map[int]map[ResponseLanguage]ErrorResponse{
	ErrCodeInvalidPhone: {
		Vietnamese: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeInvalidPhone,
			Message:     "Số điện thoại không hợp lệ",
		},
		English: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeInvalidPhone,
			Message:     "Phone number is invalid",
		},
	},
	ErrCodeInvalidName: {
		Vietnamese: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeInvalidName,
			Message:     "Tên không hợp lệ",
		},
		English: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeInvalidName,
			Message:     "Name is invalid",
		},
	},
	ErrCodeInvalidEmail: {
		Vietnamese: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeInvalidEmail,
			Message:     "Email không hợp lệ",
		},
		English: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeInvalidEmail,
			Message:     "Email is invalid",
		},
	},
	ErrCodeUserNotFound: {
		Vietnamese: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeUserNotFound,
			Message:     "Người dùng chưa tồn tại",
		},
		English: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeUserNotFound,
			Message:     "User is not exist",
		},
	},
	ErrCodeInvalidRequest: {
		Vietnamese: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeInvalidRequest,
			Message:     "Mọi thứ vẫn ổn trừ trang này, cùng dạo một vòng trước khi thử lại nhé!",
		},
		English: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeInvalidRequest,
			Message:     "Everything is fine except here, let's take a tour before trying again!",
		},
	},
	ErrCodeInternalError: {
		Vietnamese: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeInternalError,
			Message:     "Mọi thứ vẫn ổn trừ trang này, cùng dạo một vòng trước khi thử lại nhé!",
		},
		English: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeInternalError,
			Message:     "Everything is fine except here, let's take a tour before trying again!",
		},
	},
	ErrCodeUserAlreadyRequestCreateOrganization: {
		Vietnamese: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeUserAlreadyRequestCreateOrganization,
			Message:     "User đã gửi yêu cầu tạo tổ chức",
		},
		English: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeUserAlreadyRequestCreateOrganization,
			Message:     "User has already create organization request",
		},
	},
}
