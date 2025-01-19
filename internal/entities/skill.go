package entities

type Skill struct {
	Id            int    `db:"skill_id"`
	TeacherId     int    `db:"teacher_id"`
	CategoryId    int    `db:"category_id"`
	CategoryName  string `db:"category_name"`
	VideoCardLink string `db:"video_card_link"`
	About         string `db:"about"`
	Rate          int8   `db:"rate"`
	IsActive      bool   `db:"is_active"`
}
