package validator

import (
	"github.com/go-playground/validator/v10"
	"net/url"
	"regexp"
)

const urlShortNameString = "^[a-zA-Z0-9\\-]+$"

var urlShortNameRegex = regexp.MustCompile(urlShortNameString)

// ValidateUrlShortName implements validator.Func
func ValidateUrlShortName(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	if v == "" {
		return true
	}
	return urlShortNameRegex.MatchString(v)
}

// ValidateUrlScheme implements validator.Func
func ValidateUrlScheme(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	if v == "" {
		return true
	}

	u, err := url.Parse(v)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	return true
}
