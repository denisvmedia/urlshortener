package resource

import (
	"fmt"
	"github.com/denisvmedia/urlshortener/storage"
	"github.com/go-extras/api2go"
	"github.com/go-extras/errors"
	"github.com/go-playground/validator/v10"
	"net/http"
)

const resourceNotFound = "resource not found"
const validationError = "validation error"
const internalServerError = "internal server error"

var errorCodes = map[interface{}]int{
	storage.ErrNotFound:               http.StatusNotFound,
	storage.ErrShortNameAlreadyExists: http.StatusBadRequest,
}

// StatusByError gives a http error for a particular go error
func StatusByError(err error) int {
	if _, ok := err.(validator.ValidationErrors); ok {
		return http.StatusBadRequest
	}

	cause := errors.Cause(err)
	status, ok := errorCodes[cause]
	if !ok {
		status = http.StatusInternalServerError
	}
	return status
}

// HTTPErrorPtr makes an *api2go.HTTPError by the given function arguments,
// where err is an error that had occurred, msg is a text that you might want to show to the visitor,
// status is a http status you want to return along with the error message
func HTTPErrorPtr(err error, msg string, status int) *api2go.HTTPError {
	if status == 500 {
		msg = internalServerError
	}
	tmp := api2go.NewHTTPError(err, msg, status)
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			tmp.Errors = append(tmp.Errors, api2go.Error{
				ID:     fmt.Sprintf("%s:%s", e.Namespace(), e.ActualTag()),
				Title:  "Field validation failed",
				Detail: fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", e.Field(), e.ActualTag()),
				// Source: nil, // TODO: convert to "/data/attributes/*"
			})
		}
	}
	return &tmp
}

// HTTPErrorPtrWithStatus is a shortcut to HTTPErrorPtr that uses
// StatusByError for the given error to define the http status
func HTTPErrorPtrWithStatus(err error, msg string) *api2go.HTTPError {
	return HTTPErrorPtr(err, msg, StatusByError(err))
}
