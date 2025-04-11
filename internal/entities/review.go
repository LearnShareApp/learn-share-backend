package entities

type Review struct {
	ID         int    `db:"review_id"`
	TeacherID  int    `db:"teacher_id"`
	StudentID  int    `db:"student_id"`
	CategoryID int    `db:"category_id"`
	SkillID    int    `db:"skill_id"`
	Rate       int    `db:"rate"`
	Comment    string `db:"comment"`

	StudentData *User `db:"-"`
}
