package get_student_lessons

import (
	"errors"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
	"go.uber.org/zap"
	"net/http"
)

const (
	Route = "/lessons"
)

// MakeHandler returns http.HandlerFunc
// @Summary Get lessons for students
// @Description Return all lessons which have student
// @Tags students
// @Produce json
// @Success 200 {object} response
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /student/lessons [get]
// @Security     BearerAuth
func MakeHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userId := r.Context().Value(jwt.UserIDKey).(int)
		if userId == 0 {
			log.Error("id was missed in context")
			if err := httputils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		lessons, err := s.Do(r.Context(), userId)
		if err != nil {
			if errors.Is(err, serviceErrors.ErrorUserNotFound) {
				if err = httputils.RespondWith401(w, serviceErrors.ErrorUserNotFound.Error()); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return

			} else {
				log.Error(err.Error())
				if err = httputils.RespondWith500(w); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return
			}
		}

		var resp response

		if lessons != nil {
			resp = response{
				Lessons: make([]respLessons, len(lessons)),
			}

			for i := range lessons {
				resp.Lessons[i] = respLessons{
					LessonId:       lessons[i].Id,
					TeacherId:      lessons[i].TeacherId,
					TeacherUserId:  lessons[i].TeacherUserData.Id,
					TeacherEmail:   lessons[i].TeacherUserData.Email,
					TeacherName:    lessons[i].TeacherUserData.Name,
					TeacherSurname: lessons[i].TeacherUserData.Surname,
					TeacherAvatar:  lessons[i].TeacherUserData.Avatar,
					CategoryId:     lessons[i].CategoryId,
					CategoryName:   lessons[i].CategoryName,
					Status:         lessons[i].StatusName,
					Datetime:       lessons[i].ScheduleTimeDatetime,
				}
			}
		}

		respondErr := httputils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}
