package errs

import "errors"

var (
	ErrNotFoundUser      = errors.New("user not found")
	ErrIncorrectPassword = errors.New("incorrect password")
)
