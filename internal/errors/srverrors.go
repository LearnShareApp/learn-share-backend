package errors

import "errors"

var (
	ErrorPasswordTooShort  = errors.New("password too short")
	ErrorPasswordIncorrect = errors.New("password incorrect")

	ErrorUserExists       = errors.New("user already exists")
	ErrorUserNotFound     = errors.New("user not found")
	ErrorUserIsNotTeacher = errors.New("user is not a teacher")

	ErrorTeacherExists     = errors.New("teacher already exists")
	ErrorTeacherNotFound   = errors.New("teacher not found")
	ErrorSkillUnregistered = errors.New("teacher has not this skill")

	ErrorSkillRegistered = errors.New("skill already registered")

	ErrorCategoryNotFound = errors.New("category not found")

	ErrorScheduleTimeExists            = errors.New("schedule time already exists")
	ErrorScheduleTimeNotFound          = errors.New("schedule time not found")
	ErrorScheduleTimeForAnotherTeacher = errors.New("schedule time belongs to another teacher")
	ErrorScheduleTimeUnavailable       = errors.New("schedule time unavailable anymore")

	ErrorStudentAndTeacherSame     = errors.New("student and teacher the same person")
	ErrorLessonTimeBooked          = errors.New("lesson time already booked")
	ErrorLessonNotFound            = errors.New("lesson not found")
	ErrorNotRelatedUserToLesson    = errors.New("user no related to this lesson")
	ErrorNotRelatedTeacherToLesson = errors.New("teacher no related to this lesson")
	ErrorFinishedLessonNotFound    = errors.New("finished lesson not found")

	ErrorStatusNonVerification = errors.New("non-verification status")
	ErrorStatusNonWaiting      = errors.New("non-waiting status")
	ErrorStatusNonOngoing      = errors.New("non-ongoing status")

	ErrorImageNotFound       = errors.New("image not found")
	ErrorIncorrectFileFormat = errors.New("incorrect file format")

	ErrorReviewExists = errors.New("review already exists")
)
