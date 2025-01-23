package errors

import "errors"

var (
	ErrorSelectEmpty   = errors.New("select empty")
	ErrorNonUniqueData = errors.New("non unique data")
)
