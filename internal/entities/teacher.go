package entities

type Teacher struct {
	Id             int              `db:"teacher_id"`
	UserId         int              `db:"user_id"`
	Rate           float32          `db:"rate"`
	TotalRateScore int              `db:"total_rate_score"`
	ReviewsCount   int              `db:"reviews_count"`
	Skills         []*Skill         `db:"-"`
	TeacherStat    TeacherStatistic `db:"-"`
}
