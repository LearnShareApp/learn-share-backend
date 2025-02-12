package entities

type Review struct {
	Id         int    `db:"review_id"`
	TeacherId  int    `db:"teacher_id"`
	StudentId  int    `db:"student_id"`
	CategoryId int    `db:"category_id"`
	SkillId    int    `db:"skill_id"`
	Rate       int    `db:"rate"`
	Comment    string `db:"comment"`

	StudentData *User `db:"-"`
}
