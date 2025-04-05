package entities

const (
	CancelStatusName       = "cancel"
	WaitingStatusName      = "waiting"
	VerificationStatusName = "verification"
	OngoingStatusName      = "ongoing"
	FinishedStatusName     = "finished"
)

type Status struct {
	StatusId int    `db:"status_id"`
	Name     string `db:"name"`
}
