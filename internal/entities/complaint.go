package entities

import "time"

type Complaint struct {
	ID           int       `db:"complaint_id"`
	ComplainerID int       `db:"complainer_id"`
	ReportedID   int       `db:"reported_id"`
	Reason       string    `db:"reason"`
	Description  string    `db:"description"`
	CreatedAt    time.Time `db:"created_at"`

	Complainer *User `db:"-"`
	Reported   *User `db:"-"`
}
