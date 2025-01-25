package get_teacher_lessons

import (
	"errors"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/jsonutils"
	"github.com/LearnShareApp/learn-share-backend/internal/service/jwt"
	"go.uber.org/zap"
	"net/http"
)

const (
	Route = "/lessons"
)

// MakeHandler returns http.HandlerFunc
// @Summary Get lessons for teachers
// @Description Return all lessons which have teacher
// @Tags teachers
// @Produce json
// @Success 200 {object} response
// @Failure 400 {object} jsonutils.ErrorStruct
// @Failure 401 {object} jsonutils.ErrorStruct
// @Failure 403 {object} jsonutils.ErrorStruct
// @Failure 500 {object} jsonutils.ErrorStruct
// @Router /teacher/lessons [get]
// @Security     BearerAuth
func MakeHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userId := r.Context().Value(jwt.UserIDKey).(int)
		if userId == 0 {
			log.Error("id was missed in context")
			if err := jsonutils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		lessons, err := s.Do(r.Context(), userId)
		if err != nil {
			if errors.Is(err, serviceErrors.ErrorUserNotFound) {
				if err = jsonutils.RespondWith401(w, serviceErrors.ErrorUserNotFound.Error()); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return

			} else if errors.Is(err, serviceErrors.ErrorUserIsNotTeacher) {
				if err = jsonutils.RespondWith403(w, serviceErrors.ErrorUserIsNotTeacher.Error()); err != nil {
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
			Lessons: make([]respLessons, len(lessons)),
		}

		for i := range lessons {
			resp.Lessons[i] = respLessons{
				LessonId:       lessons[i].Id,
				StudentId:      lessons[i].StudentId,
				StudentName:    lessons[i].StudentUserData.Name,
				StudentSurname: lessons[i].StudentUserData.Surname,
				CategoryId:     lessons[i].CategoryId,
				Token:          lessons[i].Token,
				CategoryName:   lessons[i].CategoryName,
				Status:         lessons[i].StatusName,
				Datetime:       lessons[i].ScheduleTimeDatetime,
			}
		}

		respondErr := jsonutils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}
