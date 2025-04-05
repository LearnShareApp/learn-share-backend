package user

import (
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
)

const (
	GetPublicRoute    = "/{id}/profile"
	GetProtectedRoute = "/profile"
)

// GetUserPublic returns http.HandlerFunc
// @Summary Get user profile
// @Description Get info about user by user id in route (/api/users/{id}/profile)
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} getUserResponse
// @Failure 404 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /users/{id}/profile [get]
func (h *UserHandlers) GetUserPublic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httputils.GetIntParamFromRequestPath(r, "id")

		if err != nil {
			h.log.Error("failed to parse id from URL path", zap.Error(err))
			if err := httputils.RespondWith400(w, "missed {id} param in url path"); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}

			return
		}

		user, err := h.userService.GetUser(r.Context(), id)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				err = httputils.RespondWith404(w, err.Error())

			default:
				h.log.Error(err.Error())
				err = httputils.RespondWith500(w)
			}

			if err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}

			return
		}

		resp := mappingToUserResp(user)

		if err = httputils.SuccessRespondWith200(w, resp); err != nil {
			h.log.Error("failed to send response", zap.Error(err))
		}
	}
}

// GetUserProtected returns http.HandlerFunc
// @Summary Get user profile
// @Description Get info about user by jwt token (in Authorization enter: Bearer <your_jwt_token>)
// @Tags users
// @Produce json
// @Success 200 {object} getUserResponse
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /user/profile [get]
// @Security     BearerAuth
func (h *UserHandlers) GetUserProtected() http.HandlerFunc {
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

		user, err := h.userService.GetUser(r.Context(), id)
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

		resp := mappingToUserResp(user)

		if err = httputils.SuccessRespondWith200(w, resp); err != nil {
			h.log.Error("failed to send response", zap.Error(err))
		}
	}
}

/* helpers */

func mappingToUserResp(user *entities.User) *getUserResponse {
	resp := getUserResponse{
		ID:                  user.Id,
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

/* Mapping struct */

type getUserResponse struct {
	ID                  int       `json:"id"                   example:"1"`
	Email               string    `json:"email"                example:"qwerty@example.com"`
	Name                string    `json:"name"                 example:"John"`
	Surname             string    `json:"surname"              example:"Smith"`
	RegistrationDate    time.Time `json:"registration_date"    example:"2022-09-09T10:10:10+09:00"`
	Birthdate           time.Time `json:"birthdate"            example:"2002-09-09T10:10:10+09:00"`
	Avatar              string    `json:"avatar"               example:"uuid.png"`
	FinishedLessons     int       `json:"finished_lessons"     example:"0"`
	VerificationLessons int       `json:"verification_lessons" example:"0"`
	WaitingLessons      int       `json:"waiting_lessons"      example:"0"`
	CountOfTeachers     int       `json:"count_of_teachers"    example:"0"`
	IsTeacher           bool      `json:"is_teacher"           example:"false"`
}
