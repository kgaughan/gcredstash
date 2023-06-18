package command

import "errors"

var (
	ErrTooFewArgs  = errors.New("too few arguments")
	ErrTooManyArgs = errors.New("too many arguments")
)
