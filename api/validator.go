package api

import (
	"github.com/CodeSingerGnC/MicroBank/util"
	"github.com/go-playground/validator/v10"
)

var validCurreny validator.Func = func(fl validator.FieldLevel) bool {
	if curreny, ok := fl.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(curreny)
	}
	return false
}