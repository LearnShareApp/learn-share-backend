package get_lesson

import (
	"errors"
	"net/http"
	"strconv"

	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"go.uber.org/zap"
)

const (
	Route = "/{id}"
)

// MakeHandler returns http.HandlerFunc
// @Summary Get lesson data by lesson's id
// @Description Return lesson data by lesson's id
// @Tags lessons
// @Produce json
// @Param id path int true "LessonID"
// @Success 200 {object} response
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 404 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /lessons/{id} [get]
func MakeHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get lesson id from path
		paramId := r.PathValue("id")
		if paramId == "" {
			if err := httputils.RespondWith400(w, "missed {id} param in url path"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		lessonId, err := strconv.Atoi(paramId)
		if err != nil {
			log.Error("failed to parse id from URL path", zap.Error(err))
			if err := httputils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		lesson, err := s.Do(r.Context(), lessonId)

		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorLessonNotFound):
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

		resp := response{
			LessonId: lessonId,

			TeacherId:      lesson.TeacherId,
			TeacherUserId:  lesson.TeacherUserData.Id,
			TeacherEmail:   lesson.TeacherUserData.Email,
			TeacherName:    lesson.TeacherUserData.Name,
			TeacherSurname: lesson.TeacherUserData.Surname,
			TeacherAvatar:  lesson.TeacherUserData.Avatar,

			StudentId:      lesson.StudentId,
			StudentEmail:   lesson.StudentUserData.Email,
			StudentName:    lesson.StudentUserData.Name,
			StudentSurname: lesson.StudentUserData.Surname,
			StudentAvatar:  lesson.StudentUserData.Avatar,

			CategoryId:   lesson.CategoryId,
			CategoryName: lesson.CategoryName,
			Status:       lesson.StatusName,
			Datetime:     lesson.ScheduleTimeDatetime,
		}

		respondErr := httputils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}
