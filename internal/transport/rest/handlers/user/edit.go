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
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
)

const (
	EditRoute = "/profile"
)

// EditUser returns http.HandlerFunc
// @Summary Edit user
// @Description Edit base data about user (optional fields)
// @Tags users
// @Accept json
// @Produce json
// @Param editUserRequest body editUserRequest true "Update Info"
// @Success 200
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 413
// @Failure 500 {object} httputils.ErrorStruct
// @Router /user/profile [patch]
// @Security     BearerAuth
func (h *UserHandlers) EditUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userIDValue := r.Context().Value(jwt.UserIDKey)
		id, ok := userIDValue.(int)
		if !ok || id == 0 {
			h.log.Error("invalid or missing user ID in context", zap.Any("value", userIDValue))
			if err := httputils.RespondWith500(w); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		var req editUserRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			if err = httputils.RespondWith400(w, "failed to decode body"); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		if req.Name == "" &&
			req.Surname == "" &&
			req.Avatar == "" &&
			req.Birthdate.Before(time.Date(0001, 01, 01, 0, 0, 0, 1, time.UTC)) {
			if err := httputils.RespondWith400(w, "all fields are empty"); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		if req.Birthdate.After(time.Date(0001, 01, 01, 0, 0, 0, 1, time.UTC)) && req.Birthdate.Before(time.Date(1900, 01, 01, 0, 0, 0, 1, time.UTC)) {
			if err := httputils.RespondWith400(w, "birthdate is too old"); err != nil {
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
			Name:      req.Name,
			Surname:   req.Surname,
			Birthdate: req.Birthdate,
		}

		err := h.userService.EditUser(r.Context(), id, user, avatarReader, avatarSize)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				err = httputils.RespondWith401(w, err.Error())
			default:
				h.log.Error(err.Error())
				err = httputils.RespondWith500(w)
			}

			if err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		respondErr := httputils.SuccessRespondWith200(w, struct{}{})
		if respondErr != nil {
			h.log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

/* Mapping struct */

// @Description User registration editUserRequest.
type editUserRequest struct {
	Name      string    `json:"name"      example:"John"`
	Surname   string    `json:"surname"   example:"Smith"`
	Birthdate time.Time `json:"birthdate" example:"2000-01-01T00:00:00Z"`
	Avatar    string    `json:"avatar"    example:"base64 encoded image"`
}
