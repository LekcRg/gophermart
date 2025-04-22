package errs

import (
	"errors"
)

var (
	ErrNotFoundUser                 = errors.New("user not found")
	ErrIncorrectPassword            = errors.New("incorrect password")
	ErrAccrualReqNotRegisteredOrder = errors.New("order not registered")
	ErrOrdersRegisteredOtherUser    = errors.New("order uploaded by another user")
	ErrOrdersRegisteredThisUser     = errors.New("order uploaded by this user")
	ErrUserSmallBalance             = errors.New("small balance")
)
