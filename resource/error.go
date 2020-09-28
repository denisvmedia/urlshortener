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

func HttpErrorPtr(err error, msg string, status int) *api2go.HTTPError {
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

func HttpErrorPtrWithStatus(err error, msg string) *api2go.HTTPError {
	return HttpErrorPtr(err, msg, StatusByError(err))
}
