package registration

import "time"

type request struct {
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Password  string    `json:"password"`
	Birthdate time.Time `json:"birthdate"`
}

type response struct {
	Token string `json:"token"`
}
