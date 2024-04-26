package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslation "github.com/go-playground/validator/v10/translations/en"
)

var Translator ut.Translator

var tagValidation = make(map[string]validator.Func)

type (
	Validator interface {
		Struct(val interface{}) error
	}

	validatorIns struct {
		*validator.Validate
	}
)

func NewValidator() (Validator, error) {
	validate := validator.New()
	enTrans := en.New()
	uni := ut.New(enTrans)
	enTranslator, _ := uni.GetTranslator("en")
	if err := enTranslation.
		RegisterDefaultTranslations(validate, enTranslator); err != nil {
		return nil, err
	}
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.Split(fld.Tag.Get("json"), ",")[0]
		if name == "-" {
			return ""
		}
		return name
	})
	Translator = enTranslator
	// register the custom validation tags
	for tagName, fun := range tagValidation {
		if err := validate.RegisterValidation(tagName, fun); err != nil {
			return nil, err
		}
	}
	return &validatorIns{
		Validate: validate,
	}, nil
}

func (v *validatorIns) Struct(val interface{}) error {
	return v.Validate.Struct(val)
}

func RegisterValidation(tag string, fn validator.Func) {
	tagValidation[tag] = fn
}
