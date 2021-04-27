package models

import "strings"

const (
	ErrNotFound          = modelError("models: resource not found")
	ErrIDInvalid         = modelError("models: ID provided was invalid")
	ErrNotImplemented    = modelError("models: not implemented")
	ErrEmailRequired     = modelError("models: email address is required")
	ErrEmailInvalid      = modelError("models: email address is not valid")
	ErrEmailTaken        = modelError("models: email address is already taken")
	ErrPasswordInccorect = modelError("models: incorrect password provided")
	ErrPasswordRequired  = modelError("models: password is required")
	ErrPasswordTooShort  = modelError("models: password must be at least 8 characters long")
	ErrRememberTooShort  = modelError("models: remember token is too short")
	ErrRememberRequired  = modelError("models: remember hash is required")
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	return strings.Replace(string(e), "models: ", "", 1)
}
