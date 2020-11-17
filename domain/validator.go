package domain

import "regexp"

// define a global emailRegExp for validation (disclaimer: don't try to write this for youself!)
var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	// we want to return a map of errors for each key, pair JSON
	errors map[string]string
}

// We create a constructor for the map
func NewValidator() *Validator {
	return &Validator{errors: make(map[string]string)}
}

// We need to define a function mustbelongerthan           \/ must be bigger than this
func (v *Validator) MustBeLongerThan(field, value string, high int) bool {

	if _, ok := v.errors[field]; ok {
		return false
	}

	if value == "" {
		return true
	}

	if len(value) < high {
		v.errors[field] = ErrNotLongEnough{
			field:  field,
			amount: high,
		}.Error()

		return false
	}

	return true
}

func (v *Validator) MustBeNotEmpty(field, value string) bool {
	if _, ok := v.errors[field]; ok {
		return false
	}

	if value == "" {
		v.errors[field] = ErrIsRequired{field: field}.Error()
		return false
	}

	return true
}

func (v *Validator) MustBeValidEmail(field, email string) bool {
	if _, ok := v.errors[field]; ok {
		return false
	}

	if !emailRegexp.MatchString(email) {
		v.errors[field] = ErrEmailBadFormat.Error()
		return false
	}

	return true
}

type ElementMatcher struct {
	field string
	value string
}

func (v *Validator) MustMatch(el, match ElementMatcher) bool {
	if _, ok := v.errors[el.field]; ok {
		return false
	}

	if el.value != match.value {
		v.errors[el.field] = ErrMustMatch{match.field}.Error()
		v.errors[match.field] = ErrMustMatch{el.field}.Error()
		return false
	}

	return true
}

func (v *Validator) IsValid() bool {
	// To check if is valid, we need to return a true (no errors)
	return len(v.errors) == 0
}
