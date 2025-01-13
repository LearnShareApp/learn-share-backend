package get_profile

import (
	"github.com/LearnShareApp/learn-share-backend/internal/jsonutils"
	"time"
)

type response struct {
	Email            string    `json:"email" example:"qwerty@example.com"`
	Name             string    `json:"name" example:"John"`
	Surname          string    `json:"surname" example:"Smith"`
	RegistrationDate time.Time `json:"registration_date" example:"2022-09-09T10:10:10+09:00"`
	Birthdate        time.Time `json:"birthdate" example:"2002-09-09T10:10:10+09:00"`
}

type errorResponse jsonutils.ErrorStruct