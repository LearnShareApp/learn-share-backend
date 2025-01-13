package entities

type Category struct {
	Id     int64  `db:"category_id"`
	Name   string `db:"name"`
	MinAge int64  `db:"min_age"`
}
