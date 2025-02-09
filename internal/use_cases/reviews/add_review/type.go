package add_review

type request struct {
	TeacherId  int    `json:"teacher_id" example:"1" binding:"required"`
	CategoryId int    `json:"category_id" example:"1" binding:"required"`
	Rate       int    `json:"rate" example:"1" binding:"required"`
	Comment    string `json:"comment" example:"some comment" binding:"required"`
}
