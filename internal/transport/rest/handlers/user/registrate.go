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
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"github.com/LearnShareApp/learn-share-backend/internal/imgutils"
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
			if err = httputils.RespondWith400(w, "failed to decode body"); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		if req.Email == "" || req.Password == "" || req.Name == "" || req.Surname == "" {
			if err := httputils.RespondWith400(w, "email, name, surname or password is empty"); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		if req.Birthdate.Before(time.Date(1900, 01, 01, 0, 0, 0, 0, time.UTC)) {
			if err := httputils.RespondWith400(w, "birthdate is missed or too old"); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		var avatarReader io.Reader
		var avatarSize int64

		// if upload avatar
		if req.Avatar != "" {
			imageBytes, err := imgutils.DecodeImage(req.Avatar)

			if err != nil {
				if err = httputils.RespondWith400(w, err.Error()); err != nil {
					h.log.Error("failed to send response", zap.Error(err))
				}
				return
			}

			width, height, err := imgutils.GetImageDimensions(imageBytes)
			if err != nil {
				h.log.Error("failed to get image dimension", zap.Error(err))
				if err = httputils.RespondWith500(w); err != nil {
					h.log.Error("failed to send response", zap.Error(err))
				}
				return
			}

			if err = imgutils.CheckDimension(1, 1, width, height); err != nil {
				if err = httputils.RespondWith400(w, err.Error()); err != nil {
					h.log.Error("failed to send response", zap.Error(err))
				}
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
				err = httputils.RespondWithError(w, http.StatusConflict, err.Error())
			case errors.Is(err, serviceErrors.ErrorPasswordTooShort):
				err = httputils.RespondWith400(w, err.Error())
			default:
				h.log.Error(err.Error())
				err = httputils.RespondWith500(w)
			}

			if err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}

			return

		}

		token, err := h.jwtService.GenerateJWTToken(userID)
		if err != nil {
			h.log.Error(err.Error())
			if err = httputils.RespondWith500(w); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}

			return
		}

		var resp authResponse
		resp.Token = token

		respondErr := httputils.SuccessRespondWith201(w, resp)
		if respondErr != nil {
			h.log.Error("failed to send response", zap.Error(respondErr))
		}
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
