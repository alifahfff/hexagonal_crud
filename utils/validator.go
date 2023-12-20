package utils

import (
	validators "github.com/go-playground/validator/v10"
	"strings"
)

func customErrorMessage(fe validators.FieldError) string {
	switch fe.Tag() {
	case "required_if":
		return fe.Field() + " is required"
	case "required":
		return fe.Field() + " is required"
	case "email":
		return "Invalid email."
	case "oneof":
		return fe.Field() + " must be one of the following: " + fe.Param()
	}
	return fe.Error() // default error
}

func Validate(payload interface{}) string {
	var errMsg []string

	validate := validators.New()

	var err error

	if err = validate.Struct(payload); err == nil {
		return ""
	}

	for _, e := range err.(validators.ValidationErrors) {
		errMsg = append(errMsg, customErrorMessage(e))
	}

	return strings.Join(errMsg, ", ")
}
