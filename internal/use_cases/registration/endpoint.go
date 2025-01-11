package registration

import (
	"encoding/json"
	"errors"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/jsonutils"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const Route = "/signup"

// MakeHandler returns http.HandlerFunc
// @Summary Register new user
// @Description Register a new user (student) in the system
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request true "Registration Info"
// @Success 201 {object} response
// @Failure 400 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/signup [post]
func MakeHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			if err = jsonutils.RespondWith400(w, "failed to decode body"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		if req.Email == "" || req.Password == "" || req.Name == "" || req.Surname == "" { //  || req.Avatar == ""
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

		user := &entities.User{
			Email:     req.Email,
			Password:  req.Password,
			Name:      req.Name,
			Surname:   req.Surname,
			Birthdate: req.Birthdate,
		}

		token, err := s.Do(r.Context(), user)
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
