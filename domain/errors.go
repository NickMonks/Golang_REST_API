package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNoResult                     = errors.New("no result")
	ErrUserWithEmailAlreadyExist    = errors.New("user with email already exist")
	ErrUserWithUsernameAlreadyExist = errors.New("user with username already exist")
	ErrEmailBadFormat               = errors.New("Error: Email not valid")
	ErrInvalidCredential            = errors.New("Error: Invalid credentials")
)

type ErrNotLongEnough struct {
	field  string
	amount int
}

func (e ErrNotLongEnough) Error() string {
	return fmt.Sprintf("%v not long enough; %d characters is required", e.field, e.amount)
}

// we want this errnotlongmenough be a error type. to do so, because error is an interface
// we must implement the function

type ErrIsRequired struct {
	field string
}

func (e ErrIsRequired) Error() string {
	return fmt.Sprintf("%v is required", e.field)
}

type ErrMustMatch struct {
	field string
}

func (e ErrMustMatch) Error() string {
	return fmt.Sprintf("must match %v", e.field)
}
