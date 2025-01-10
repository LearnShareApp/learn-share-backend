package entities

type Category struct {
	Id     int64  `db:"id"`
	Name   string `db:"name"`
	MinAge int64  `db:"min_age"`
}
