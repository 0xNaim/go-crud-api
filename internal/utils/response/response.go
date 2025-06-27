package response

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

const (
	StatusOk    = "Ok"
	StatusError = "Error"
)

func WriteJson(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status:  StatusError,
		Message: err.Error(),
	}

}

func ValidateError(errs validator.ValidationErrors) Response {
	var errMessages []string

	for _, err := range errs {
		field := err.Field()
		param := err.Param()

		switch err.ActualTag() {
		case "required":
			errMessages = append(errMessages, field+" is required.")
		case "min":
			errMessages = append(errMessages, field+" must be at least "+param+" characters long.")
		case "max":
			errMessages = append(errMessages, field+" must be at most "+param+" characters long.")
		case "email":
			errMessages = append(errMessages, field+" must be a valid email address.")
		case "numeric":
			errMessages = append(errMessages, field+" must be a numeric value.")
		default:
			errMessages = append(errMessages, field+" is invalid.")
		}
	}

	return Response{
		Status:  StatusError,
		Message: strings.Join(errMessages, " "),
	}
}
