package entities

const (
	CancelStatusName       = "cancelled"
	WaitingStatusName      = "waiting"
	VerificationStatusName = "verification"
	OngoingStatusName      = "ongoing"
	FinishedStatusName     = "finished"
)

type Status struct {
	StatusId int    `db:"status_id"`
	Name     string `db:"name"`
}
