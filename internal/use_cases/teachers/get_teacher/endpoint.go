package get_teacher

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
	ProtectedRoute = ""
	PublicRoute    = "/{id}"
)

// MakeProtectedHandler returns http.HandlerFunc which handle get teacher, get user id from token
// @Summary Get teacher data
// @Description Get all info about teacher (user info + teacher + his skills) by user id in token
// @Tags teachers
// @Produce json
// @Success 200 {object} response
// @Failure 403 {object} httputils.ErrorStruct
// @Failure 404 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /teacher [get]
// @Security     BearerAuth
func MakeProtectedHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userId := r.Context().Value(jwt.UserIDKey).(int)
		if userId == 0 {
			log.Error("id was missed in context")
			if err := httputils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		user, err := s.DoByUserId(r.Context(), userId)

		if err != nil {
			coveringErrors(w, log, err)
			return
		}

		resp := mappingToResponse(user)

		respondErr := httputils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

// MakePublicHandler returns http.HandlerFunc which handle get teacher, get user id from http param
// @Summary Get teacher data
// @Description Get all info about teacher (user info + teacher + his skills) by his TeacherID in route (/api/teachers/{id})
// @Tags teachers
// @Produce json
// @Param id path int true "Teacher's ID"
// @Success 200 {object} response
// @Failure 404 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /teachers/{id} [get]
func MakePublicHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teacherId, err := httputils.GetIntParamFromRequestPath(r, "id")
		if err != nil {
			if err := httputils.RespondWith400(w, "missed or not-number {id} param in url path"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		user, err := s.DoByTacherId(r.Context(), teacherId)

		if err != nil {
			coveringErrors(w, log, err)
			return
		}

		resp := mappingToResponse(user)

		respondErr := httputils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

func coveringErrors(w http.ResponseWriter, log *zap.Logger, err error) {
	switch {
	case errors.Is(err, serviceErrors.ErrorUserIsNotTeacher):
		err = httputils.RespondWith403(w, serviceErrors.ErrorUserIsNotTeacher.Error())
	case errors.Is(err, serviceErrors.ErrorTeacherNotFound):
		err = httputils.RespondWith404(w, serviceErrors.ErrorTeacherNotFound.Error())
	default:
		log.Error(err.Error())
		err = httputils.RespondWith500(w)
	}

	if err != nil {
		log.Error("failed to send response", zap.Error(err))
	}
}

func mappingToResponse(user *entities.User) *response {
	resp := response{
		TeacherId:          user.TeacherData.Id,
		UserId:             user.Id,
		Email:              user.Email,
		Name:               user.Name,
		Surname:            user.Surname,
		RegistrationDate:   user.RegistrationDate,
		Birthdate:          user.Birthdate,
		Avatar:             user.Avatar,
		FinishedLessons:    user.TeacherData.TeacherStat.CountOfFinishedLesson,
		CountOfStudents:    user.TeacherData.TeacherStat.CountOfStudents,
		CommonRate:         user.TeacherData.Rate,
		CommonReviewsCount: user.TeacherData.ReviewsCount,

		Skills: make([]skill, 0, len(user.TeacherData.Skills)),
	}

	// remap entity skill to response skill-type
	for _, sk := range user.TeacherData.Skills {
		resp.Skills = append(resp.Skills, skill{
			SkillId:       sk.Id,
			CategoryId:    sk.CategoryId,
			CategoryName:  sk.CategoryName,
			VideoCardLink: sk.VideoCardLink,
			About:         sk.About,
			Rate:          sk.Rate,
			ReviewsCount:  sk.ReviewsCount,
		})
	}

	return &resp
}
