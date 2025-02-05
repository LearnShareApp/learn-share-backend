package entities

type Teacher struct {
	Id          int              `db:"teacher_id"`
	UserId      int              `db:"user_id"`
	Skills      []*Skill         `db:"-"`
	TeacherStat TeacherStatistic `db:"-"`
}
