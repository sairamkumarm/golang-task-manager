package utils

import "github.com/go-playground/validator/v10"

var Validate = validator.New()

func ValidateStruct(v interface{}) error {
    return Validate.Struct(v)
}
