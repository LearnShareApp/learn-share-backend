package entities

import "time"

type ScheduleTime struct {
	Id          int       `db:"schedule_time_id"`
	TeacherId   int       `db:"teacher_id"`
	Datetime    time.Time `db:"datetime"`
	IsAvailable bool      `db:"is_available"`
}
