package sso

type SessionData struct {
	IdentityID    string
	IsActive      bool
	EmailVerified bool
	Email         string
}
