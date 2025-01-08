package login

type request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type response struct {
	Token string `json:"token"`
}
