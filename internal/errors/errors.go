package errors

import "errors"

var (
	ErrorUserExists      = errors.New("user already exists")
	ErrorUserNotFound    = errors.New("user not found")
	ErrorTeacherExists   = errors.New("teacher already exists")
	ErrorSkillRegistered = errors.New("skill already registered")

	ErrorPasswordTooShort  = errors.New("password too short")
	ErrorPasswordIncorrect = errors.New("password incorrect")

	ErrorCategoryNotFound = errors.New("category not found")
)
