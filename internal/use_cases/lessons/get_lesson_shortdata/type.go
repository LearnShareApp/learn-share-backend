package get_lesson_shortdata

type response struct {
	LessonId      int    `json:"lesson_id" example:"1"`
	TeacherId     int    `json:"teacher_id" example:"1"`
	TeacherUserId int    `json:"teacher_user_id" example:"1"`
	StudentId     int    `json:"student_id" example:"1"`
	CategoryId    int    `json:"category_id" example:"1"`
	CategoryName  string `json:"category_name" example:"Programming"`
}
