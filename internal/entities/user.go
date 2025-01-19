package entities

import "time"

type User struct {
	Id               int       `db:"user_id"`
	Email            string    `db:"email"`
	Name             string    `db:"name"`
	Surname          string    `db:"surname"`
	Password         string    `db:"password"`
	RegistrationDate time.Time `db:"registration_date"`
	Birthdate        time.Time `db:"birthdate"`
	Avatar           string    `db:"avatar"`
	IsTeacher        bool      `db:"-"`
	TeacherData      *Teacher  `db:"-"`
}
