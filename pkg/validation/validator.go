package validation

import (
	"sync"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	vLib "github.com/go-playground/validator/v10"
	ru_translations "github.com/go-playground/validator/v10/translations/ru"
)

type Validator struct {
	v     *vLib.Validate
	trans ut.Translator
}

// Field describes a field in a struct.
// It is used to get info about it for validation.
type Field vLib.FieldLevel

var (
	singleton *Validator
	once      sync.Once
)

func GetValidator() *Validator {
	once.Do(func() {
		uni := ut.New(en.New())
		trans, _ := uni.GetTranslator("ru")
		v := vLib.New()
		ru_translations.RegisterDefaultTranslations(v, trans)

		singleton = &Validator{
			v:     v,
			trans: trans,
		}
	})

	return singleton
}

// RegisterValidation registers a new validation function with the validator.
// This function will be called when the validator encounters the tag.
func (v Validator) RegisterValidation(tag string, check func(Field) bool, errMsg string) error {
	err := v.v.RegisterTranslation(tag, v.trans,
		func(ut ut.Translator) error {
			return ut.Add(tag, errMsg, true)
		},
		func(ut ut.Translator, fe vLib.FieldError) string {
			t, err := ut.T(tag, fe.Field())
			if err != nil {
				// TODO: change this to something more error tolerant
				panic("could not register validation")
			}
			return t
		},
	)
	if err != nil {
		return err
	}
	return v.v.RegisterValidation(tag, func(fl vLib.FieldLevel) bool {
		return check(fl)
	})
}

// ValidateStruct validates the given struct.
// It validates it according to the tags on the struct.
// You can lookup available tags [here].
// Or create custom tags by using RegisterValidation.
//
// [here]: https://github.com/go-playground/validator
func (v Validator) Validate(value any) error {
	return v.v.Struct(value)
}

// UnpackErrors unpacks the error returned by ValidateStruct into a slice of strings.
func (v Validator) UnpackErrors(e error) []string {
	values, ok := e.(vLib.ValidationErrors)
	if !ok {
		return nil
	}
	errs := make([]string, 0, len(values))
	for _, vv := range values {
		errs = append(errs, vv.Translate(v.trans))
	}
	return errs
}

func (v Validator) Var(field string, tag string) error {
	return v.v.Var(field, tag)
}
