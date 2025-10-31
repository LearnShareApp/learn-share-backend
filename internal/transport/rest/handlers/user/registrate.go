package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/imgutils"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/httputils"
)

const RegistrationRoute = "/signup"

// RegistrationUser returns http.HandlerFunc
// @Summary Register new user
// @Description Register a new user (student) in the system
// @Tags auth
// @Accept json
// @Produce json
// @Param registrationRequest body registrationRequest true "Registration Info"
// @Success 201 {object} authResponse
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 409 {object} httputils.ErrorStruct
// @Failure 413
// @Failure 500 {object} httputils.ErrorStruct
// @Router /auth/signup [post]
func (h *UserHandlers) RegistrationUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req registrationRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httputils.RespondWith400(w, "failed to decode body", h.log)

			return
		}

		if req.Email == "" || req.Password == "" || req.Name == "" || req.Surname == "" {
			httputils.RespondWith400(w, "email, name, surname or password is empty", h.log)

			return
		}

		if req.Birthdate.Before(time.Date(1920, 01, 01, 0, 0, 0, 0, time.UTC)) {
			httputils.RespondWith400(w, "birthdate is missed or too old", h.log)

			return
		}

		var avatarReader io.Reader
		var avatarSize int64

		// if upload avatar
		if req.Avatar != "" {
			imageBytes, err := imgutils.DecodeImage(req.Avatar)

			if err != nil {
				httputils.RespondWith400(w, err.Error(), h.log)

				return
			}

			width, height, err := imgutils.GetImageDimensions(imageBytes)
			if err != nil {
				h.log.Error("failed to get image dimension", zap.Error(err))
				httputils.RespondWith500(w, h.log)

				return
			}

			if err = imgutils.CheckDimension(1, 1, width, height); err != nil {
				httputils.RespondWith400(w, err.Error(), h.log)

				return
			}
			avatarReader = bytes.NewReader(imageBytes)
			defer func() {
				if closer, ok := avatarReader.(io.Closer); ok {
					if err := closer.Close(); err != nil {
						h.log.Error("failed to close reader", zap.Error(err))
					}
				}
			}()
			avatarSize = int64(len(imageBytes))
		}

		user := &entities.User{
			Email:     req.Email,
			Password:  req.Password,
			Name:      req.Name,
			Surname:   req.Surname,
			Birthdate: req.Birthdate,
		}

		userID, err := h.userService.CreateUser(r.Context(), user, avatarReader, avatarSize)

		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserExists):
				httputils.RespondWith409(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorPasswordTooShort):
				httputils.RespondWith400(w, err.Error(), h.log)
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

		httputils.SuccessRespondWith201(w, resp, h.log)
	}
}

// @Description User registration registrationRequest.
type registrationRequest struct {
	Email     string    `json:"email"     example:"john@gmail.com"       binding:"required,email"`
	Name      string    `json:"name"      example:"John"                 binding:"required"`
	Surname   string    `json:"surname"   example:"Smith"                binding:"required"`
	Password  string    `json:"password"  example:"strongpass123"        binding:"required"`
	Birthdate time.Time `json:"birthdate" example:"2000-01-01T00:00:00Z" binding:"required"`
	Avatar    string    `json:"avatar"    example:"base64 encoded image"`
}
