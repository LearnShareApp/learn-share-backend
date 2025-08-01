package entities

type StateMachineName string

const (
	LessonStateMachineName StateMachineName = "lesson"
	SkillStateMachineName  StateMachineName = "skill"
)

type StateMachineItem struct {
	ID             int    `db:"item_id"`
	StateMachineID int    `db:"state_machine_id"`
	StateID        int    `db:"state_id"`
	StateName      string `db:"state_name"`
}
