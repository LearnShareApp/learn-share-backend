package errors

import "errors"

var (
	ErrorUserExists        = errors.New("user already exists")
	ErrorUserNotFound      = errors.New("user not found")
	ErrorPasswordTooShort  = errors.New("password too short")
	ErrorPasswordIncorrect = errors.New("password incorrect")
)
