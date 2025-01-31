package get_user

import (
	"errors"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/jsonutils"
	"github.com/LearnShareApp/learn-share-backend/internal/service/jwt"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

const (
	PublicRoute    = "/{id}/profile"
	ProtectedRoute = "/profile"
)

// MakeProtectedHandler returns http.HandlerFunc
// @Summary Get user profile
// @Description Get info about user by jwt token (in Authorization enter: Bearer <your_jwt_token>)
// @Tags users
// @Produce json
// @Success 200 {object} response
// @Failure 401 {object} jsonutils.ErrorStruct
// @Failure 500 {object} jsonutils.ErrorStruct
// @Router /user/profile [get]
// @Security     BearerAuth
func MakeProtectedHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.Context().Value(jwt.UserIDKey).(int)
		if id == 0 {
			log.Error("id was missed in context")
			if err := jsonutils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		user, err := s.Do(r.Context(), id)
		if err != nil {
			if errors.Is(err, serviceErrors.ErrorUserNotFound) {
				if err = jsonutils.RespondWith401(w, err.Error()); err != nil {
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

		resp := response{
			Id:               user.Id,
			Email:            user.Email,
			Name:             user.Name,
			Surname:          user.Surname,
			RegistrationDate: user.RegistrationDate,
			Birthdate:        user.Birthdate,
			Avatar:           user.Avatar,
			IsTeacher:        user.IsTeacher,
		}

		if err = jsonutils.SuccessRespondWith200(w, resp); err != nil {
			log.Error("failed to send response", zap.Error(err))
		}
	}
}

// MakePublicHandler returns http.HandlerFunc
// @Summary Get user profile
// @Description Get info about user by user id in route (/api/users/{id}/profile)
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response
// @Failure 404 {object} jsonutils.ErrorStruct
// @Failure 500 {object} jsonutils.ErrorStruct
// @Router /users/{id}/profile [get]
func MakePublicHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id int

		paramId := r.PathValue("id")
		if paramId == "" {
			if err := jsonutils.RespondWith400(w, "missed {id} param in url path"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		id, err := strconv.Atoi(paramId)

		if err != nil {
			log.Error("failed to parse id from URL path", zap.Error(err))
			if err := jsonutils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		user, err := s.Do(r.Context(), id)
		if err != nil {
			if errors.Is(err, serviceErrors.ErrorUserNotFound) {
				if err := jsonutils.RespondWith404(w, serviceErrors.ErrorUserNotFound.Error()); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return
			}

			log.Error(err.Error())
			if err = jsonutils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		resp := response{
			Id:               user.Id,
			Email:            user.Email,
			Name:             user.Name,
			Surname:          user.Surname,
			RegistrationDate: user.RegistrationDate,
			Birthdate:        user.Birthdate,
			IsTeacher:        user.IsTeacher,
		}

		if err = jsonutils.SuccessRespondWith200(w, resp); err != nil {
			log.Error("failed to send response", zap.Error(err))
		}
	}
}
