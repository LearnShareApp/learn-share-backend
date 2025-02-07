package get_lesson_shortdata

import (
	"errors"
	"net/http"
	"strconv"

	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/jsonutils"
	"go.uber.org/zap"
)

const (
	Route = "/{id}/short-data"
)

// MakeHandler returns http.HandlerFunc
// @Summary Get lesson really short data by lesson's id
// @Description Return lesson short data by lesson's id
// @Tags lessons
// @Produce json
// @Param id path int true "LessonID"
// @Success 200 {object} response
// @Failure 400 {object} jsonutils.ErrorStruct
// @Failure 404 {object} jsonutils.ErrorStruct
// @Failure 500 {object} jsonutils.ErrorStruct
// @Router /lessons/{id}/short-data [get]
func MakeHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get lesson id from path
		paramId := r.PathValue("id")
		if paramId == "" {
			if err := jsonutils.RespondWith400(w, "missed {id} param in url path"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		lessonId, err := strconv.Atoi(paramId)
		if err != nil {
			log.Error("failed to parse id from URL path", zap.Error(err))
			if err := jsonutils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		lesson, err := s.Do(r.Context(), lessonId)

		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorLessonNotFound):
				err = jsonutils.RespondWith404(w, err.Error())
			default:
				log.Error(err.Error())
				err = jsonutils.RespondWith500(w)
			}

			if err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		resp := response{
			LessonId: lessonId,

			TeacherId:     lesson.TeacherId,
			TeacherUserId: lesson.TeacherUserData.Id,

			StudentId: lesson.StudentId,

			CategoryId:   lesson.CategoryId,
			CategoryName: lesson.CategoryName,
		}

		respondErr := jsonutils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}
