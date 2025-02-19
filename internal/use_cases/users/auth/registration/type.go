package registration

import (
	"time"
)

// @Description User registration request
type request struct {
	Email     string    `json:"email" example:"john@gmail.com" binding:"required,email"`
	Name      string    `json:"name" example:"John" binding:"required"`
	Surname   string    `json:"surname" example:"Smith" binding:"required"`
	Password  string    `json:"password" example:"strongpass123" binding:"required"`
	Birthdate time.Time `json:"birthdate" example:"2000-01-01T00:00:00Z" binding:"required"`
	Avatar    string    `json:"avatar" example:"base64 encoded image"`
}

// @Description User registration response
type response struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
