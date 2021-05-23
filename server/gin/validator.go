package gin

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	universalTranslator "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/spf13/cast"
)

type ValidateFuncs struct {

	TagName string
	Fn validator.Func
	Message string
}


var Trans universalTranslator.Translator

func InitValidator(validateFuncs []*ValidateFuncs) error {
	var err error
	var ok bool
	var v *validator.Validate
	if v, ok = binding.Validator.Engine().(*validator.Validate); ok {
		for _, f := range validateFuncs{
			fmt.Printf("RegisterValidation func:%s\n", f.TagName)
			err = v.RegisterValidation(f.TagName, f.Fn)
			if err!=nil{
				return err
			}
		}
	}

	return initTrans(validateFuncs)
}

// initTrans ...
func initTrans(validateFuns []*ValidateFuncs) (err error) {
	local := "en"
	var ok bool
	var v *validator.Validate
	if v, ok = binding.Validator.Engine().(*validator.Validate); ok {
		enT := en.New()
		uni := universalTranslator.New(enT, enT)
		Trans, ok = uni.GetTranslator(local)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", local)
		}
		err = enTranslations.RegisterDefaultTranslations(v, Trans)
		if err != nil {
			return err
		}
		for _, f := range validateFuns{
			fmt.Printf("RegisterTranslation func:%s message:%s\n", f.TagName, f.Message)
			err = v.RegisterTranslation(
				f.TagName,
				Trans,
				registerTranslator(f.TagName, f.Message),
				translate,
			)
			if err != nil {
				return err
			}
		}

		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			return fld.Tag.Get("json")
		})


	}
	return nil
}

// registerTranslator ...
func registerTranslator(tag string, msg string) validator.RegisterTranslationsFunc {
	return func(trans universalTranslator.Translator) error {
		if err := trans.Add(tag, msg, false); err != nil {
			return err
		}
		return nil
	}
}

// translate ....
func translate(trans universalTranslator.Translator, fe validator.FieldError) string {
	msg, err := trans.T(fe.Tag(), fe.Field(), fe.Param())
	if err != nil {
		panic(fe.(error).Error())
	}
	return msg
}

// Refer to here to write a custom check function
// https://github.com/go-playground/validator/blob/f6584a41c8acc5dfc0b62f7962811f5231c11530/baked_in.go
// https://github.com/go-playground/validator/issues/524

func IsLess(fl validator.FieldLevel) bool {

	value := fl.Field().Int()
	param:= fl.Param()
	return value < cast.ToInt64(param)
}