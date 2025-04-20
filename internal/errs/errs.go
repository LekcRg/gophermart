package errs

import (
	"errors"
)

var (
	ErrNotFoundUser              = errors.New("user not found")
	ErrIncorrectPassword         = errors.New("incorrect password")
	ErrOrdersRegisteredOtherUser = errors.New(
		"order has already been uploaded by another user")
	ErrOrdersRegisteredThisUser = errors.New(
		"order has already been uploaded by this user")
	ErrUserSmallBalance = errors.New("small balance")
)
