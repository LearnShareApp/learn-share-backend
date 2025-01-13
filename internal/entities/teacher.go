package entities

type Teacher struct {
	Id     int64 `db:"teacher_id"`
	UserId int64 `db:"user_id"`
}
