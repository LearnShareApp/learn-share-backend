package entities

import "time"

type Lesson struct {
	Id             int `db:"lesson_id"`
	StudentId      int `db:"student_id"`
	TeacherId      int `db:"teacher_id"`
	CategoryId     int `db:"category_id"`
	ScheduleTimeId int `db:"schedule_time_id"`
	StatusId       int `db:"status_id"`
	Price          int `db:"price"`

	StatusName           string    `db:"-"`
	CategoryName         string    `db:"-"`
	ScheduleTimeDatetime time.Time `db:"-"`
	StudentUserData      *User     `db:"-"` // info about student
	TeacherUserData      *User     `db:"-"` // info about teacher (as user)
}
