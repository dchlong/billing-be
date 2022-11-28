package presenter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error *Error `json:"error"`
}

type Error struct {
	Code       string `json:"code"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func NewInvalidUserNameError(maxLength int) *ErrorResponse {
	err := &Error{
		Code:       "invalid_user_name",
		StatusCode: http.StatusBadRequest,
		Message:    fmt.Sprintf("username must be from 1-%d characters", maxLength),
	}

	return &ErrorResponse{
		Error: err,
	}
}

func NewInvalidInputError(inputObj interface{}, err error) *ErrorResponse {
	message := err.Error()
	uerr, ok := err.(*json.UnmarshalTypeError)
	if ok {
		message = fmt.Sprintf("field %s must be not %s", uerr.Field, uerr.Value)
	}

	verrs, ok := err.(validator.ValidationErrors)
	if ok {
		verr := verrs[0]
		fieldName := verr.Field()
		field, _ := reflect.TypeOf(inputObj).Elem().FieldByName(fieldName)
		fieldJSONName, _ := field.Tag.Lookup("json")
		if verr.Tag() == "gt" {
			message = fmt.Sprintf("field %s must be greater than %s", fieldJSONName, verr.Param())
		} else if verr.Tag() == "required" {
			message = fmt.Sprintf("%s is a required field", fieldJSONName)
		}
	}
	newerr := &Error{
		Code:       "invalid_input",
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}

	return &ErrorResponse{
		Error: newerr,
	}
}

func NewInternalServerErrorError(err error) *ErrorResponse {
	newErr := &Error{
		Code:       "internal_server_error",
		StatusCode: http.StatusInternalServerError,
		Message:    err.Error(),
	}

	return &ErrorResponse{
		Error: newErr,
	}
}
