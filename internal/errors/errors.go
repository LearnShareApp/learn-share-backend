package errors

import "errors"

var (
	ErrorUserExists       = errors.New("user already exists")
	ErrorUserNotFound     = errors.New("user not found")
	ErrorTeacherExists    = errors.New("teacher already exists")
	ErrorTeacherNotFound  = errors.New("teacher not found")
	ErrorUserIsNotTeacher = errors.New("user is not teacher")
	ErrorSkillRegistered  = errors.New("skill already registered")

	ErrorPasswordTooShort  = errors.New("password too short")
	ErrorPasswordIncorrect = errors.New("password incorrect")

	ErrorCategoryNotFound   = errors.New("category not found")
	ErrorScheduleTimeExists = errors.New("schedule time already exists")

	// repository
	ErrorSelectEmpty = errors.New("select empty")
)
