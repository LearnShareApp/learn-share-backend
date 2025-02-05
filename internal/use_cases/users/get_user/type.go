package get_user

import (
	"time"
)

type response struct {
	Id                  int       `json:"id" example:"1"`
	Email               string    `json:"email" example:"qwerty@example.com"`
	Name                string    `json:"name" example:"John"`
	Surname             string    `json:"surname" example:"Smith"`
	RegistrationDate    time.Time `json:"registration_date" example:"2022-09-09T10:10:10+09:00"`
	Birthdate           time.Time `json:"birthdate" example:"2002-09-09T10:10:10+09:00"`
	Avatar              string    `json:"avatar" example:"uuid.png"`
	FinishedLessons     int       `json:"finished_lessons" example:"0"`
	VerificationLessons int       `json:"verification_lessons" example:"0"`
	WaitingLessons      int       `json:"waiting_lessons" example:"0"`
	CountOfTeachers     int       `json:"count_of_teachers" example:"0"`
	IsTeacher           bool      `json:"is_teacher" example:"false"`
}
