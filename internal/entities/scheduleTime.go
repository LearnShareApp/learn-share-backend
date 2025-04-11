package entities

import "time"

type ScheduleTime struct {
	ID          int       `db:"schedule_time_id"`
	TeacherID   int       `db:"teacher_id"`
	Datetime    time.Time `db:"datetime"`
	IsAvailable bool      `db:"is_available"`
}
