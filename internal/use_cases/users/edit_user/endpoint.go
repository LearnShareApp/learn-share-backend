package edit_user

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/imgutils"
	"github.com/LearnShareApp/learn-share-backend/internal/jsonutils"
	"github.com/LearnShareApp/learn-share-backend/internal/service/jwt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

const Route = "/profile"

// MakeHandler returns http.HandlerFunc
// @Summary Edit user
// @Description Edit base data about user (optional fields)
// @Tags users
// @Accept json
// @Produce json
// @Param request body request true "Update Info"
// @Success 200
// @Failure 400 {object} jsonutils.ErrorStruct
// @Failure 401 {object} jsonutils.ErrorStruct
// @Failure 413
// @Failure 500 {object} jsonutils.ErrorStruct
// @Router /user/profile [patch]
// @Security     BearerAuth
func MakeHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.Context().Value(jwt.UserIDKey).(int)
		if id == 0 {
			log.Error("id was missed in context")
			if err := jsonutils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		var req request

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			if err = jsonutils.RespondWith400(w, "failed to decode body"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		if req.Password == "" &&
			req.Name == "" &&
			req.Surname == "" &&
			req.Avatar == "" &&
			req.Birthdate.Before(time.Date(0001, 01, 01, 0, 0, 0, 1, time.UTC)) {
			if err := jsonutils.RespondWith400(w, "all fields are empty"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		if req.Birthdate.After(time.Date(0001, 01, 01, 0, 0, 0, 1, time.UTC)) && req.Birthdate.Before(time.Date(1900, 01, 01, 0, 0, 0, 1, time.UTC)) {
			if err := jsonutils.RespondWith400(w, "birthdate is too old"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		var avatarReader io.Reader
		var avatarSize int64

		// if upload avatar
		if req.Avatar != "" {
			imageBytes, err := imgutils.DecodeImage(req.Avatar)

			if err != nil {
				if err = jsonutils.RespondWith400(w, err.Error()); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return
			}

			width, height, err := imgutils.GetImageDimensions(imageBytes)
			if err != nil {
				log.Error("failed to get image dimension", zap.Error(err))
				if err = jsonutils.RespondWith500(w); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return
			}

			if err = imgutils.CheckDimension(1, 1, width, height); err != nil {
				if err = jsonutils.RespondWith400(w, err.Error()); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return
			}
			avatarReader = bytes.NewReader(imageBytes)
			defer func() {
				if closer, ok := avatarReader.(io.Closer); ok {
					if err := closer.Close(); err != nil {
						log.Error("failed to close reader", zap.Error(err))
					}
				}
			}()
			avatarSize = int64(len(imageBytes))
		}

		user := &entities.User{
			Password:  req.Password,
			Name:      req.Name,
			Surname:   req.Surname,
			Birthdate: req.Birthdate,
		}

		err := s.Do(r.Context(), id, user, avatarReader, avatarSize)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				err = jsonutils.RespondWith401(w, err.Error())
			case errors.Is(err, serviceErrors.ErrorPasswordTooShort):
				err = jsonutils.RespondWith400(w, err.Error())
			default:
				log.Error(err.Error())
				err = jsonutils.RespondWith500(w)
			}

			if err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		respondErr := jsonutils.SuccessRespondWith200(w, struct{}{})
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}
