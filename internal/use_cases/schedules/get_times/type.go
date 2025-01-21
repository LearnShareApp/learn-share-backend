package get_times

import "time"

type response struct {
	Datetimes []time.Time `json:"datetimes" example:"0001-01-01T00:00:00Z"`
}
