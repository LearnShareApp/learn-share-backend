package book_lesson

// @Description book lesson body request
type request struct {
	TeacherId      int `json:"teacher_id" example:"1" binding:"required"` // @Description exactly teacherID, not his userID
	CategoryId     int `json:"category_id" example:"1" binding:"required"`
	ScheduleTimeId int `json:"schedule_time_id" example:"1" binding:"required"`
}

// response isn't needed, just 201 if success
