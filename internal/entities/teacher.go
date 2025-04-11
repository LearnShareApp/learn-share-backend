package entities

type Teacher struct {
	ID             int              `db:"teacher_id"`
	UserID         int              `db:"user_id"`
	Rate           float32          `db:"rate"`
	TotalRateScore int              `db:"total_rate_score"`
	ReviewsCount   int              `db:"reviews_count"`
	Skills         []*Skill         `db:"-"`
	TeacherStat    TeacherStatistic `db:"-"`
}
