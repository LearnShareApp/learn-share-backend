package registration

type request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type response struct {
	JwtToken string `json:"jwt_token"`
}
