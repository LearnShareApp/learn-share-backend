package entities

type Skill struct {
	Id             int     `db:"skill_id"`
	TeacherId      int     `db:"teacher_id"`
	CategoryId     int     `db:"category_id"`
	CategoryName   string  `db:"category_name"`
	VideoCardLink  string  `db:"video_card_link"`
	About          string  `db:"about"`
	Rate           float32 `db:"rate"`
	TotalRateScore int     `db:"total_rate_score"`
	ReviewsCount   int     `db:"reviews_count"`
	IsActive       bool    `db:"is_active"`
}
