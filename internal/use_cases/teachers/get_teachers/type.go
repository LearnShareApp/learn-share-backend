package get_teachers

import "time"

type response struct {
	Teachers []teacher `json:"teachers"`
}

type teacher struct {
	TeacherId        int       `json:"teacher_id" example:"1"`
	UserId           int       `json:"user_id" example:"1"`
	Email            string    `json:"email" example:"qwerty@example.com"`
	Name             string    `json:"name" example:"John"`
	Surname          string    `json:"surname" example:"Smith"`
	RegistrationDate time.Time `json:"registration_date" example:"2022-09-09T10:10:10+09:00"`
	Birthdate        time.Time `json:"birthdate" example:"2002-09-09T10:10:10+09:00"`
	Avatar           string    `json:"avatar" example:"uuid.png"`
	Skills           []skill   `json:"skills"`
}

type skill struct {
	SkillId       int    `json:"skill_id" example:"1"`
	CategoryId    int    `json:"category_id" example:"1"`
	CategoryName  string `json:"category_name" example:"Category"`
	VideoCardLink string `json:"video_card_link" example:"https://youtu.be/HIcSWuKMwOw?si=FtxN1QJU9ZWnXy85"`
	About         string `json:"about" example:"about me..."`
	Rate          int8   `json:"rate" example:"5"`
}
