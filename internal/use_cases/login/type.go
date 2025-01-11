package login

import "github.com/LearnShareApp/learn-share-backend/internal/jsonutils"

type request struct {
	Email    string `json:"email" example:"john@gmail.com" binding:"required,email"`
	Password string `json:"password" example:"strongpass123" binding:"required"`
}

type response struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type errorResponse jsonutils.ErrorStruct
