package entities

type Statuses struct {
	StatusId int    `db:"status_id"`
	Name     string `db:"name"`
}
