package models

import "strings"

const (
	ErrNotFound          = modelError("models: resource not found")
	ErrEmailRequired     = modelError("models: email address is required")
	ErrEmailInvalid      = modelError("models: email address is not valid")
	ErrEmailTaken        = modelError("models: email address is already taken")
	ErrPasswordInccorect = modelError("models: incorrect password provided")
	ErrPasswordRequired  = modelError("models: password is required")
	ErrPasswordTooShort  = modelError("models: password must be at least 8 characters long")
	ErrTitleRequired     = modelError("models: title is required")

	ErrNotImplemented   = privateError("models: not implemented")
	ErrRememberTooShort = privateError("models: remember token is too short")
	ErrIDInvalid        = privateError("models: ID provided was invalid")
	ErrRememberRequired = privateError("models: remember hash is required")
	ErrUserIDRequired   = privateError("models: user ID is required")
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	return strings.Replace(string(e), "models: ", "", 1)
}

type privateError string

func (e privateError) Error() string {
	return string(e)
}
