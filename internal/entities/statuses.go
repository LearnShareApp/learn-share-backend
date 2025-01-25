package entities

const (
	CancelStatusName = "cancelled"
)

type Status struct {
	StatusId int    `db:"status_id"`
	Name     string `db:"name"`
}
