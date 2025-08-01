package entities

type StateMachine struct {
	ID           int    `db:"state_machine_id"`
	Name         string `db:"name"`
	StartStateID int    `db:"start_state_id"`
}
