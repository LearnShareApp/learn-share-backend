package get_teacher_lessons

import "time"

type response struct {
	Lessons []respLessons `json:"lessons"`
}

type respLessons struct {
	LessonId       int       `json:"lesson_id" example:"1"`
	StudentId      int       `json:"student_id" example:"1"`
	StudentName    string    `json:"student_name" example:"John"`
	StudentSurname string    `json:"student_surname" example:"Smith"`
	CategoryId     int       `json:"category_id" example:"1"`
	Token          string    `json:"token" example:"headerrrrrrrrrrrrr.bbbbbbbodybbbbbbbbbbb.subscribeeeeeeeeee"`
	CategoryName   string    `json:"category_name" example:"Programming"`
	Status         string    `json:"status" example:"verification"`
	Datetime       time.Time `json:"datetime" example:"2025-02-01T09:00:00Z"`
}
