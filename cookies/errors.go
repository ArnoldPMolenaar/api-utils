package cookies

import "errors"

var (
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")
	ErrInvalidConfiguration = errors.New("invalid cookie configuration")
)
