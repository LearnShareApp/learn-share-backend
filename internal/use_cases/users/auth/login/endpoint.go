package login

import (
	"encoding/json"
	"errors"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"go.uber.org/zap"
	"net/http"
)

const Route = "/login"

// MakeHandler returns http.HandlerFunc
// @Summary Login user
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request true "Login Credentials"
// @Success 200 {object} response
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /auth/login [post]
func MakeHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			if err = httputils.RespondWith400(w, "failed to decode body"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		if req.Email == "" || req.Password == "" {
			if err := httputils.RespondWith400(w, "email or password is empty"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		user := &entities.User{
			Email:    req.Email,
			Password: req.Password,
		}

		token, err := s.Do(r.Context(), user)
		if err != nil {
			if errors.Is(err, serviceErrors.ErrorUserNotFound) {
				if err = httputils.RespondWith401(w, err.Error()); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return

			} else if errors.Is(err, serviceErrors.ErrorPasswordIncorrect) {
				if err = httputils.RespondWith401(w, err.Error()); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return

			} else {
				log.Error(err.Error())
				if err = httputils.RespondWith500(w); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return
			}
		}

		var resp response
		resp.Token = token

		respondErr := httputils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}
