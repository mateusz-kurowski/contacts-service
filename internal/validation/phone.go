package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/phonenumbers"
)

func PhoneNumberValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		num, ok := fl.Field().Interface().(string)
		if !ok || num == "" {
			return false
		}

		phoneNum, err := phonenumbers.Parse(num, "PL")
		if err != nil {
			return false
		}

		if !phonenumbers.IsPossibleNumber(phoneNum) {
			return false
		}

		return phonenumbers.IsValidNumber(phoneNum)
	}
}
