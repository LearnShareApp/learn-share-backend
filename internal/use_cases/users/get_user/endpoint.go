package get_user

import (
	"errors"

	"net/http"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"github.com/LearnShareApp/learn-share-backend/internal/service/jwt"

	"go.uber.org/zap"
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
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /user/profile [get]
// @Security     BearerAuth
func MakeProtectedHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(jwt.UserIDKey).(int)
		if id == 0 {
			log.Error("id was missed in context")
			if err := httputils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		user, err := s.Do(r.Context(), id)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				err = httputils.RespondWith401(w, err.Error())

			default:
				log.Error(err.Error())
				err = httputils.RespondWith500(w)
			}

			if err != nil {
				log.Error("failed to send response", zap.Error(err))
			}

			return
		}

		resp := mappingToResponse(user)

		if err = httputils.SuccessRespondWith200(w, resp); err != nil {
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
// @Failure 404 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /users/{id}/profile [get]
func MakePublicHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httputils.GetIntParamFromRequestPath(r, "id")

		if err != nil {
			log.Error("failed to parse id from URL path", zap.Error(err))
			if err := httputils.RespondWith400(w, "missed {id} param in url path"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}

			return
		}

		user, err := s.Do(r.Context(), id)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				err = httputils.RespondWith404(w, err.Error())

			default:
				log.Error(err.Error())
				err = httputils.RespondWith500(w)
			}

			if err != nil {
				log.Error("failed to send response", zap.Error(err))
			}

			return
		}

		resp := mappingToResponse(user)

		if err = httputils.SuccessRespondWith200(w, resp); err != nil {
			log.Error("failed to send response", zap.Error(err))
		}
	}
}

func mappingToResponse(user *entities.User) *response {
	resp := response{
		Id:                  user.Id,
		Email:               user.Email,
		Name:                user.Name,
		Surname:             user.Surname,
		RegistrationDate:    user.RegistrationDate,
		Birthdate:           user.Birthdate,
		Avatar:              user.Avatar,
		FinishedLessons:     user.Stat.CountOfFinishedLesson,
		VerificationLessons: user.Stat.CountOfVerificationLesson,
		WaitingLessons:      user.Stat.CountOfWaitingLesson,
		CountOfTeachers:     user.Stat.CountOfTeachers,
		IsTeacher:           user.IsTeacher,
	}

	return &resp
}
