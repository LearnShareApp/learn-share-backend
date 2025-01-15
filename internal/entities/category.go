package entities

type Category struct {
	Id     int    `db:"category_id"`
	Name   string `db:"name"`
	MinAge int    `db:"min_age"`
}
