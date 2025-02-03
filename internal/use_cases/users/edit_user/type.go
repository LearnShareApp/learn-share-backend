package edit_user

import (
	"time"
)

// @Description User registration request
type request struct {
	Name      string    `json:"name" example:"John"`
	Surname   string    `json:"surname" example:"Smith"`
	Password  string    `json:"password" example:"strongpass123"`
	Birthdate time.Time `json:"birthdate" example:"2000-01-01T00:00:00Z"`
	Avatar    string    `json:"avatar" example:"base64 encoded image"`
}
