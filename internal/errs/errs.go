package errs

import "errors"

var (
	NotFoundUser      = errors.New("user not found")
	IncorrectPassword = errors.New("incorrect password")
)
