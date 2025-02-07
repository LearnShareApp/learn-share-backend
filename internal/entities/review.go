package entities

type Review struct {
	ReviewId  int    `db:"review_id"`
	TeacherId int    `db:"teacher_id"`
	StudentId int    `db:"student_id"`
	SkillId   int    `db:"skill_id"`
	Rate      int    `db:"rate"`
	Comment   string `db:"comment"`
}
