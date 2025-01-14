package get_profile

import (
	"errors"
	errors2 "github.com/LearnShareApp/learn-share-backend/internal/errors"
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
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /users/profile [get]
// @Security     BearerAuth
func MakeProtectedHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.Context().Value(jwt.UserIDKey).(int64)
		if id == 0 {
			log.Error("id was missed in context")
			if err := jsonutils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		user, err := s.Do(r.Context(), id)
		if err != nil {
			// Вроде как нет смысла обрабатывать кейс когда пользователь не найден по id т. к. по хорошему в токене,
			// который выпустили мы не может быть несуществующий пользователь
			// что и есть 500-я

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
		}

		respondErr := jsonutils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

// MakePublicHandler returns http.HandlerFunc
// @Summary Get user profile
// @Description Get info about user by id in route (/api/users/{id}/profile)
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /users/{id}/profile [get]
func MakePublicHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id int64

		paramId := r.PathValue("id")
		if paramId == "" {
			if err := jsonutils.RespondWith400(w, "missed {id} param in url path"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		id, err := strconv.ParseInt(paramId, 10, 64)

		if err != nil {
			log.Error("failed to parse id from URL path", zap.Error(err))
			if err := jsonutils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		user, err := s.Do(r.Context(), id)
		if err != nil {
			if errors.Is(err, errors2.ErrorUserNotFound) {
				log.Error("failed to find user", zap.Int64("id", id), zap.Error(err))
				if err := jsonutils.RespondWith404(w, errors2.ErrorUserNotFound.Error()); err != nil {
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
		}

		respondErr := jsonutils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}
