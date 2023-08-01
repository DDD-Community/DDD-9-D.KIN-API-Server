package user

import (
	"d.kin-app/internal/typex"
	"d.kin-app/internal/validation"
	"github.com/go-playground/validator/v10"
	"reflect"
)

func init() {
	validation.Default.RegisterValidation("gender", validationGender)
}

func validationGender(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() != reflect.String {
		return false
	}
	return typex.Contains(GenderValues(), Gender(field.String()))
}
