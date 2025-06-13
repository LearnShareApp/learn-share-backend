package entities

type StateName string

const (
	Pending    StateName = "pending"
	Approved   StateName = "approved"
	Rejected   StateName = "rejected"
	Planned    StateName = "planned"
	Cancel     StateName = "cancel"
	Ongoing    StateName = "ongoing"
	Finished   StateName = "finished"
	Conflicted StateName = "conflicted"
	Completed  StateName = "completed"
)

type State struct {
	ID   int       `db:"state_id"`
	Name StateName `db:"name"`
}
