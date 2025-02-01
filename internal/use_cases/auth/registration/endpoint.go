package registration

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/imgutils"
	"github.com/LearnShareApp/learn-share-backend/internal/jsonutils"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
	"time"
)

const Route = "/signup"
const MaxImageWeight = 5 << 20

// MakeHandler returns http.HandlerFunc
// @Summary Register new user
// @Description Register a new user (student) in the system
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request true "Registration Info"
// @Success 201 {object} response
// @Failure 400 {object} jsonutils.ErrorStruct
// @Failure 409 {object} jsonutils.ErrorStruct
// @Failure 413
// @Failure 500 {object} jsonutils.ErrorStruct
// @Router /auth/signup [post]
func MakeHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			if err = jsonutils.RespondWith400(w, "failed to decode body"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		if req.Email == "" || req.Password == "" || req.Name == "" || req.Surname == "" {
			if err := jsonutils.RespondWith400(w, "email, name, surname or password is empty"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		if req.Birthdate.Before(time.Date(1900, 01, 01, 0, 0, 0, 0, time.UTC)) {
			if err := jsonutils.RespondWith400(w, "birthdate is missed or too old"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		var avatarReader io.Reader
		var avatarSize int64

		// if upload avatar
		if req.Avatar != "" {
			avatarBytes, err := base64.StdEncoding.DecodeString(req.Avatar)
			if err != nil {
				if err = jsonutils.RespondWith400(w, "invalid avatar format"); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return
			}

			// check for weight
			if len(avatarBytes) > MaxImageWeight {
				if err = jsonutils.RespondWith400(w, "avatar is too large"); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return
			}

			// check is MIME-type (image)
			mimeType := http.DetectContentType(avatarBytes)
			if !strings.HasPrefix(mimeType, "image/") {
				if err = jsonutils.RespondWith400(w, "file is not an image"); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return
			}

			width, height, err := imgutils.GetImageDimensions(avatarBytes)
			if err != nil {
				log.Error("failed to get image dimension", zap.Error(err))
				if err = jsonutils.RespondWith500(w); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return
			}

			if width != height {
				if err = jsonutils.RespondWith400(w, "avatar must have 1:1 aspect ratio"); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return
			}

			avatarReader = bytes.NewReader(avatarBytes)
			avatarSize = int64(len(avatarBytes))
		}

		user := &entities.User{
			Email:     req.Email,
			Password:  req.Password,
			Name:      req.Name,
			Surname:   req.Surname,
			Birthdate: req.Birthdate,
		}

		token, err := s.Do(r.Context(), user, avatarReader, avatarSize)

		if err != nil {
			if errors.Is(err, serviceErrors.ErrorUserExists) {
				if err = jsonutils.RespondWithError(w, http.StatusConflict, err.Error()); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return

			} else if errors.Is(err, serviceErrors.ErrorPasswordTooShort) {
				if err = jsonutils.RespondWith400(w, err.Error()); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return

			} else {
				log.Error(err.Error())
				if err = jsonutils.RespondWith500(w); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return
			}
		}

		var resp response
		resp.Token = token

		respondErr := jsonutils.SuccessRespondWith201(w, resp)
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}
