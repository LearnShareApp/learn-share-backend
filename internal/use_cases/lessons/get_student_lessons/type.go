package get_student_lessons

import "time"

type response struct {
	Lessons []respLessons `json:"lessons"`
}

type respLessons struct {
	LessonId       int       `json:"lesson_id" example:"1"`
	TeacherId      int       `json:"teacher_id" example:"1"`
	TeacherUserId  int       `json:"teacher_user_id" example:"1"`
	TeacherName    string    `json:"teacher_name" example:"John"`
	TeacherSurname string    `json:"teacher_surname" example:"Smith"`
	CategoryId     int       `json:"category_id" example:"1"`
	CategoryName   string    `json:"category_name" example:"Programming"`
	Status         string    `json:"status" example:"verification"`
	Datetime       time.Time `json:"datetime" example:"2025-02-01T09:00:00Z"`
}
