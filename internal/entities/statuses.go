package entities

const (
	CancelStatusName       = "cancelled"
	WaitingStatusName      = "waiting"
	VerificationStatusName = "verification"
	OngoingStatusName      = "ongoing"
)

type Status struct {
	StatusId int    `db:"status_id"`
	Name     string `db:"name"`
}
