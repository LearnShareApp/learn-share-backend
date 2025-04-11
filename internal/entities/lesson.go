package entities

import "time"

type Lesson struct {
	ID             int `db:"lesson_id"`
	StudentID      int `db:"student_id"`
	TeacherID      int `db:"teacher_id"`
	CategoryID     int `db:"category_id"`
	ScheduleTimeID int `db:"schedule_time_id"`
	StatusID       int `db:"status_id"`
	Price          int `db:"price"`

	StatusName           string    `db:"-"`
	CategoryName         string    `db:"-"`
	ScheduleTimeDatetime time.Time `db:"-"`
	StudentUserData      *User     `db:"-"` // info about student
	TeacherUserData      *User     `db:"-"` // info about teacher (as user)
}
