package get_teacher

import (
	"errors"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/jsonutils"
	"github.com/LearnShareApp/learn-share-backend/internal/service/jwt"
	"go.uber.org/zap"
	"net/http"
	"strconv"
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
// @Failure 401 {object} jsonutils.ErrorStruct
// @Failure 404 {object} jsonutils.ErrorStruct
// @Failure 500 {object} jsonutils.ErrorStruct
// @Router /teacher [get]
// @Security     BearerAuth
func MakeProtectedHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userId := r.Context().Value(jwt.UserIDKey).(int)
		if userId == 0 {
			log.Error("id was missed in context")
			if err := jsonutils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		user, err := s.Do(r.Context(), userId)

		if err != nil {
			coveringErrors(w, log, err)
			return
		}

		resp := mappingToResponse(user)

		respondErr := jsonutils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

// MakePublicHandler returns http.HandlerFunc which handle get teacher, get user id from http param
// @Summary Get teacher data
// @Description Get all info about teacher (user info + teacher + his skills) by his UserID in route (/api/teachers/{id})
// @Tags teachers
// @Produce json
// @Param id path int true "Teacher's UserID"
// @Success 200 {object} response
// @Failure 404 {object} jsonutils.ErrorStruct
// @Failure 500 {object} jsonutils.ErrorStruct
// @Router /teachers/{id} [get]
func MakePublicHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var teacherUserId int

		paramId := r.PathValue("id")
		if paramId == "" {
			if err := jsonutils.RespondWith400(w, "missed {id} param in url path"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		teacherUserId, err := strconv.Atoi(paramId)

		if err != nil {
			log.Error("failed to parse id from URL path", zap.Error(err))
			if err := jsonutils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		user, err := s.Do(r.Context(), teacherUserId)

		if err != nil {
			coveringErrors(w, log, err)
			return
		}

		resp := mappingToResponse(user)

		respondErr := jsonutils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

func coveringErrors(w http.ResponseWriter, log *zap.Logger, err error) {
	if errors.Is(err, serviceErrors.ErrorUserNotFound) {
		if err := jsonutils.RespondWith404(w, serviceErrors.ErrorUserNotFound.Error()); err != nil {
			log.Error("failed to send response", zap.Error(err))
		}
	} else if errors.Is(err, serviceErrors.ErrorTeacherNotFound) {
		if err := jsonutils.RespondWith404(w, serviceErrors.ErrorUserIsNotTeacher.Error()); err != nil {
			log.Error("failed to send response", zap.Error(err))
		}
	} else {
		log.Error(err.Error())
		if err = jsonutils.RespondWith500(w); err != nil {
			log.Error("failed to send response", zap.Error(err))
		}
	}
}

func mappingToResponse(user *entities.User) *response {
	resp := response{
		TeacherId:        user.TeacherData.Id,
		UserId:           user.Id,
		Email:            user.Email,
		Name:             user.Name,
		Surname:          user.Surname,
		RegistrationDate: user.RegistrationDate,
		Birthdate:        user.Birthdate,
		Avatar:           user.Avatar,
		FinishedLessons:  user.TeacherData.TeacherStat.CountOfFinishedLesson,
		CountOfStudents:  user.TeacherData.TeacherStat.CountOfStudents,
		Skills:           make([]skill, 0, len(user.TeacherData.Skills)),
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
		})
	}

	return &resp
}
