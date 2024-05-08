package validator

import (
	"be-capstone-project/src/cmd/public/config"
	"errors"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
	"sync"
)

// TagMessage tag message for struct validate
const TagMessage = "message"

// TagErrorCode tag code for struct validation in return if the fields got error
const TagErrorCode = "errorCode"

var once sync.Once
var instance *CustomValidate

type CustomValidate struct {
	*validator.Validate
}

func (customValidate *CustomValidate) init(validate *validator.Validate) {
	customValidate.Validate = validate
}

// InitCustomValidator custom validator init and register global validator for validate request
func InitCustomValidator(cfg *config.Config) *CustomValidate {
	TSValidate := validator.New()
	_ = TSValidate.RegisterValidation("PhoneValidator", PhoneValidatorFunc)
	_ = TSValidate.RegisterValidation("EmailValidator", EmailValidatorFunc)
	customValidate := &CustomValidate{}
	customValidate.init(TSValidate)
	instance = customValidate
	return customValidate
}

func getInstance() *CustomValidate {
	//once.Do(func() {
	//	fmt.Println("once do")
	//	validate := validator.New()
	//	instance = &CustomValidate{
	//		Validate: validate,
	//	}
	//})
	return instance
}

func Struct(obj interface{}) error {
	errValidate := getInstance().Validate.Struct(obj)
	if errValidate != nil {
		var outErrorMsg strings.Builder
		for _, err := range errValidate.(validator.ValidationErrors) {
			t := reflect.TypeOf(obj)

			for i := 0; i < t.NumField(); i++ {
				if t.Field(i).Name == err.Field() {
					errMsg := t.Field(i).Tag.Get(TagMessage)
					outErrorMsg.WriteString(errMsg)
					outErrorMsg.WriteString("\n ")
				}
			}
		}
		return errors.New(outErrorMsg.String())
	}
	return nil
}
