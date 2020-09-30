package validator

import (
	"github.com/go-playground/validator/v10"
	"net/url"
	"regexp"
)

const urlShortNameString = "^[a-zA-Z0-9\\-]+$"
const blackListedValuesString = "^(api|swagger|metrics)$"

var urlShortNameRegex = regexp.MustCompile(urlShortNameString)
var blackListedValuesRegex = regexp.MustCompile(blackListedValuesString)

// ValidateUrlShortName implements validator.Func
func ValidateUrlShortName(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	if v == "" {
		return true
	}

	if blackListedValuesRegex.MatchString(v) {
		// blacklisted value
		return false
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
