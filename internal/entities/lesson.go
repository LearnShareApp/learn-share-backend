package entities

import "time"

type Lesson struct {
	Id             int    `db:"lesson_id"`
	StudentId      int    `db:"student_id"`
	TeacherId      int    `db:"teacher_id"`
	CategoryId     int    `db:"category_id"`
	ScheduleTimeId int    `db:"schedule_time_id"`
	StatusId       int    `db:"status_id"`
	Price          int    `db:"price"`
	Token          string `db:"token"`

	StatusName           string    `db:"-"`
	CategoryName         string    `db:"-"`
	ScheduleTimeDatetime time.Time `db:"-"`
	StudentData          *User     `db:"-"`
	TeacherData          *Teacher  `db:"-"`
}
