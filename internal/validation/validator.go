package validation

import (
	"github.com/go-playground/validator/v10"
	"time"
)

func init() {
	Default.RegisterValidation("year_of_birth", validationYearOfBirth)
}

func validationYearOfBirth(fl validator.FieldLevel) bool {
	field := fl.Field()
	var v int64
	if field.CanInt() {
		v = field.Int()
	} else if field.CanUint() {
		v = int64(field.Uint())
	}

	return v >= 1900 &&
		v <= int64(time.Now().Year()-9) // TODO: 한국 나이 기준 몇살 부터 가입 가능 하게 해야할까
}

var Default = validator.New()

func Valid(data any) error {
	var wrapper struct {
		Data any `validate:"dive"`
	}

	wrapper.Data = data
	return Default.Struct(&wrapper)
}
