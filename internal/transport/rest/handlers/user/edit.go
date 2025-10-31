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
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
)

const (
	editRoute = "/profile"
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
			httputils.RespondWith500(w, h.log)

			return
		}

		var req editUserRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httputils.RespondWith400(w, "failed to decode body", h.log)

			return
		}

		if req.Name == "" {
			httputils.RespondWith400(w, "field \"name\" is required", h.log)
		}

		if req.Surname == "" {
			httputils.RespondWith400(w, "field \"surname\" is required", h.log)
		}

		if req.Birthdate.Before(time.Date(1920, 01, 01, 0, 0, 0, 1, time.UTC)) {
			httputils.RespondWith400(w, "birth date is too old or empty", h.log)

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
			Name:      req.Name,
			Surname:   req.Surname,
			Birthdate: req.Birthdate,
		}

		err := h.userService.EditUser(r.Context(), id, user, avatarReader, avatarSize)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				httputils.RespondWith401(w, err.Error(), h.log)
			default:
				h.log.Error(err.Error())
				httputils.RespondWith500(w, h.log)
			}

			return
		}

		httputils.SuccessRespondWith200(w, struct{}{}, h.log)
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
