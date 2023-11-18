package helpers

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var (
	ErrInternal   = errors.New("it's our fault, not yours")
	ErrNotAllowed = errors.New("you're not allowed to perform this action")
)

type ResponseError struct {
	err  error
	code int
}

func NewResponseError(err error, code int) ResponseError {
	return ResponseError{
		err:  err,
		code: code,
	}
}

func (e ResponseError) Error() string {
	return e.err.Error()
}

func (e ResponseError) Code() int {
	return e.code
}

type validationError struct {
	Field   string
	Message string
}

func (v validationError) String() string {
	return fmt.Sprintf("%s : %s", v.Field, v.Message)
}

func GetValidationError(errValidation validator.ValidationErrors) []string {
	errs := make([]string, len(errValidation))

	for i, e := range errValidation {
		vErr := validationError{
			Field: e.Field(),
		}
		switch e.Tag() {
		case "required":
			vErr.Message = "can't be empty"
		case "email":
			vErr.Message = "must be a valid e-mail (ex: johndoe@mail.com)"
		case "min":
			vErr.Message = fmt.Sprintf("must be at least %s characters long", e.Param())
		case "url":
			vErr.Message = "must be a valid URL (ex: http://example.org)"
		default:
			vErr.Message = e.Error()
		}

		errs[i] = vErr.String()
	}

	return errs
}
