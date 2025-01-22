package get_times

import "time"

type response struct {
	Datetimes []times `json:"datetimes"`
}

type times struct {
	ScheduleTimeId int       `json:"schedule_time_id" example:"1"`
	Datetime       time.Time `json:"datetime" example:"0001-01-01T00:00:00Z"`
	IsAvailable    bool      `json:"is_available" example:"true"`
}
