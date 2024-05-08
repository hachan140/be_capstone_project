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

const (
	ErrCodeInvalidRequest = 400001

	ErrCodeChannelIsMissing   = 400012
	ErrCodeInternalError      = 500001
	ErrCodeOTPMaxRequest      = 400003
	ErrCodeOTPWrongThreeTimes = 400004
	ErrCodeOTPWrongMaxTimes   = 400005
	ErrCodeOTPResend          = 400006
	MessageInputOTPOverNTimes = 400007

	ErrCodeUserNotFound        = 400011
	ErrCodeInvalidEmail        = 400010
	ErrCodeUserIsDeletedOnApp  = 400061
	ErrCodeInvalidPhone        = 400002
	ErrCodeUserWaitingToDelete = 666409
	ErrCodeUserIsInactive      = 400027
	ErrCodeRecordExisted       = 400409
	ErrCodeChannelNotMatch     = 400009

	ErrorCodeWrapAccountBlock    = 401666 // wrap mã lỗi 666400
	ErrorCodeCannotLogin         = 666402 // Không thể login tài khoản này
	ErrorCodeCannotCreateAccount = 666401 // Không thể tạo thêm tài khoản
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
	ErrCodeUserIsDeletedOnApp: {
		Vietnamese: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeUserIsDeletedOnApp,
			Message:     "Tài khoản đã bị xóa",
		},
		English: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeUserIsDeletedOnApp,
			Message:     "This account has been deleted",
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
	//ErrorCodeWrapAccountBlock: {
	//	Vietnamese: {
	//		HTTPCode:    http.StatusBadRequest,
	//		ServiceCode: ErrorCodeWrapAccountBlock,
	//		Message:     "Mọi thứ vẫn ổn trừ trang này, cùng dạo một vòng trước khi thử lại nhé!",
	//	},
	//	English: {
	//		HTTPCode:    http.StatusBadRequest,
	//		ServiceCode: ErrorCodeWrapAccountBlock,
	//		Message:     "Everything is fine except here, let's take a tour before trying again!",
	//	},
	//},
	//ErrorCodeCannotLogin: {
	//	Vietnamese: {
	//		HTTPCode:    http.StatusBadRequest,
	//		ServiceCode: ErrorCodeCannotLogin,
	//		Message:     "Mọi thứ vẫn ổn trừ trang này, cùng dạo một vòng trước khi thử lại nhé!",
	//	},
	//	English: {
	//		HTTPCode:    http.StatusBadRequest,
	//		ServiceCode: ErrorCodeCannotLogin,
	//		Message:     "Everything is fine except here, let's take a tour before trying again!",
	//	},
	//},
	//ErrorCodeCannotCreateAccount: {
	//	Vietnamese: {
	//		HTTPCode:    http.StatusBadRequest,
	//		ServiceCode: ErrorCodeCannotCreateAccount,
	//		Message:     "Mọi thứ vẫn ổn trừ trang này, cùng dạo một vòng trước khi thử lại nhé!",
	//	},
	//	English: {
	//		HTTPCode:    http.StatusBadRequest,
	//		ServiceCode: ErrorCodeCannotCreateAccount,
	//		Message:     "Everything is fine except here, let's take a tour before trying again!",
	//	},
	//},
	ErrCodeChannelNotMatch: {
		Vietnamese: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeChannelNotMatch,
			Message:     "Mọi thứ vẫn ổn trừ trang này, cùng dạo một vòng trước khi thử lại nhé!",
		},
		English: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeChannelNotMatch,
			Message:     "Everything is fine except here, let's take a tour before trying again!",
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

	ErrCodeChannelIsMissing: {
		Vietnamese: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeChannelIsMissing,
			Message:     "Mọi thứ vẫn ổn trừ trang này, cùng dạo một vòng trước khi thử lại nhé!",
		},
		English: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeChannelIsMissing,
			Message:     "Everything is fine except here, let's take a tour before trying again!",
		},
	},
	ErrCodeOTPMaxRequest: {
		Vietnamese: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeOTPMaxRequest,
			Message:     "Bạn đã yêu cầu gửi OTP quá số lần tối đa. Vui lòng thử lại sau 24 giờ.",
		},
		English: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeOTPMaxRequest,
			Message:     "You have exceeded the maximum number of OTP request attempts. Please try again after 24 hours",
		},
	},
	ErrCodeOTPWrongThreeTimes: {
		Vietnamese: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeOTPWrongThreeTimes,
			Message:     "Bạn đã yêu cầu gửi OTP quá số lần tối đa. Vui lòng thử lại sau 24 giờ.",
		},
		English: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeOTPWrongThreeTimes,
			Message:     "You have exceeded the maximum number of OTP request attempts. Please try again after 24 hours",
		},
	},
	ErrCodeOTPWrongMaxTimes: {
		Vietnamese: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeOTPWrongMaxTimes,
			Message:     "Bạn đã yêu cầu gửi OTP quá số lần tối đa. Vui lòng thử lại sau 24 giờ.",
		},
		English: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeOTPWrongMaxTimes,
			Message:     "You have exceeded the maximum number of OTP request attempts. Please try again after 24 hours",
		},
	},
	ErrCodeOTPResend: {
		Vietnamese: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeOTPResend,
			Message:     "Vui lòng đợi %v giây để yêu cầu gửi lại mã OTP",
		},
		English: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeOTPResend,
			Message:     "Please wait %v seconds to request OTP resend",
		},
	},
	MessageInputOTPOverNTimes: {
		Vietnamese: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: MessageInputOTPOverNTimes,
			Message:     "Bạn đã nhập sai %v lần liên tiếp, vui lòng đợi %v giây để tiếp tục",
		},
		English: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: MessageInputOTPOverNTimes,
			Message:     "You entered incorrect OTP %v times, try again after %vs",
		},
	},
	ErrCodeUserWaitingToDelete: {
		Vietnamese: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeUserWaitingToDelete,
			Message:     "Tài khoản này đang được chờ để xoá",
		},
		English: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeUserWaitingToDelete,
			Message:     "This account is waiting to be deleted",
		},
	},
	ErrCodeUserIsInactive: {
		Vietnamese: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeUserIsInactive,
			Message:     "Tài khoản không hợp lệ. Bạn vui lòng liên hệ CSKH để được hỗ trợ.",
		},
		English: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeUserIsInactive,
			Message:     "This account is invalid. Please contact CS for more support.",
		},
	},
	ErrCodeRecordExisted: {
		Vietnamese: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeRecordExisted,
			Message:     "Người dùng đã tồn tại",
		},
		English: {
			HTTPCode:    http.StatusBadRequest,
			ServiceCode: ErrCodeRecordExisted,
			Message:     "User is already exist",
		},
	},
}
