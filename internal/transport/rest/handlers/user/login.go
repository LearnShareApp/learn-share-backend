package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/httputils"
)

const LoginRoute = "/login"

// LoginUser returns http.HandlerFunc
// @Summary Login user
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body loginRequest true "Login Credentials"
// @Success 200 {object} authResponse
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /auth/login [post]
func (h *UserHandlers) LoginUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httputils.RespondWith400(w, "failed to decode body", h.log)

			return
		}

		if req.Email == "" || req.Password == "" {
			httputils.RespondWith400(w, "email or password is empty", h.log)

			return
		}

		user := &entities.User{
			Email:    req.Email,
			Password: req.Password,
		}

		userID, err := h.userService.CheckUser(r.Context(), user)

		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				httputils.RespondWith401(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorPasswordIncorrect):
				httputils.RespondWith401(w, err.Error(), h.log)
			default:
				h.log.Error(err.Error())
				httputils.RespondWith500(w, h.log)
			}

			return
		}

		token, err := h.jwtService.GenerateJWTToken(userID)
		if err != nil {
			h.log.Error(err.Error())
			httputils.RespondWith500(w, h.log)

			return
		}

		var resp authResponse
		resp.Token = token

		httputils.SuccessRespondWith200(w, resp, h.log)
	}
}

type loginRequest struct {
	Email    string `json:"email"    example:"john@gmail.com" binding:"required,email"`
	Password string `json:"password" example:"strongpass123"  binding:"required"`
}

// @Description User registration authResponse.
type authResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
