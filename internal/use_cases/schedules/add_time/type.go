package add_time

import "time"

type request struct {
	Datetime time.Time `json:"datetime" example:"2025-02-01T00:00:00Z" binding:"required"`
}
