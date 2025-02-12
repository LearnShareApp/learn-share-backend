package get_reviews

type response struct {
	Reviews []review `json:"reviews"`
}

type review struct {
	ReviewId   int    `json:"review_id" example:"1"`
	TeacherId  int    `json:"teacher_id" example:"1"`
	SkillId    int    `json:"skill_id" example:"1"`
	CategoryId int    `json:"category_id" example:"1"`
	Rate       int    `json:"rate" example:"5"`
	Comment    string `json:"comment" example:"This is a comment"`

	StudentId      int    `json:"student_id" example:"1"`
	StudentEmail   string `json:"student_email" example:"qwerty@example.com"`
	StudentName    string `json:"student_name" example:"John"`
	StudentSurname string `json:"student_surname" example:"Smith"`
	StudentAvatar  string `json:"student_avatar" example:"uuid.png"`
}
