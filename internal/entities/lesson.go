package entities

type Lesson struct {
	Id             int    `db:"lesson_id"`
	StudentId      int    `db:"student_id"`
	TeacherId      int    `db:"teacher_id"`
	CategoryId     int    `db:"category_id"`
	ScheduleTimeId int    `db:"schedule_time_id"`
	StatusId       int    `db:"status_id"`
	Price          int    `db:"price"`
	Token          string `db:"token"`
}
