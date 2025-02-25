package add_skill

type request struct {
	CategoryID    int    `json:"category_id"     example:"1"                                                binding:"required"`
	VideoCardLink string `json:"video_card_link" example:"https://youtu.be/HIcSWuKMwOw?si=FtxN1QJU9ZWnXy85"`
	About         string `json:"about"           example:"I am Groot"`
}
