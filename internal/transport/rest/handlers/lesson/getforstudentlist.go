package lesson

import (
	"errors"
	"net/http"
	"time"

	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
	"go.uber.org/zap"
)

const (
	getForStudentListRoute = "/student/lessons"
)

// GetForStudentList returns http.HandlerFunc
// @Summary Get lessons for students
// @Description Return all lessons which have student
// @Tags students
// @Produce json
// @Success 200 {object} getStudentLessonsResponse
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /student/lessons [get]
// @Security     BearerAuth
func (h *LessonHandlers) GetForStudentList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDValue := r.Context().Value(jwt.UserIDKey)
		userID, ok := userIDValue.(int)
		if !ok || userID == 0 {
			h.log.Error("invalid or missing user ID in context", zap.Any("value", userIDValue))
			if err := httputils.RespondWith500(w); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		lessons, err := h.lessonService.GetStudentLessonList(r.Context(), userID)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				err = httputils.RespondWith401(w, serviceErrors.ErrorUserNotFound.Error())
			default:
				h.log.Error(err.Error())
				err = httputils.RespondWith500(w)
			}

			if err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}

			return
		}

		var resp getStudentLessonsResponse

		if lessons != nil {
			resp = getStudentLessonsResponse{
				Lessons: make([]respStudentLessons, len(lessons)),
			}

			for i := range lessons {
				resp.Lessons[i] = respStudentLessons{
					LessonID:       lessons[i].ID,
					TeacherID:      lessons[i].TeacherID,
					TeacherUserID:  lessons[i].TeacherUserData.ID,
					TeacherEmail:   lessons[i].TeacherUserData.Email,
					TeacherName:    lessons[i].TeacherUserData.Name,
					TeacherSurname: lessons[i].TeacherUserData.Surname,
					TeacherAvatar:  lessons[i].TeacherUserData.Avatar,
					CategoryID:     lessons[i].CategoryID,
					CategoryName:   lessons[i].CategoryName,
					Status:         lessons[i].StatusName,
					Datetime:       lessons[i].ScheduleTimeDatetime,
				}
			}
		}

		respondErr := httputils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			h.log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

type getStudentLessonsResponse struct {
	Lessons []respStudentLessons `json:"lessons"`
}

type respStudentLessons struct {
	LessonID       int       `json:"lesson_id"       example:"1"`
	TeacherID      int       `json:"teacher_id"      example:"1"`
	TeacherUserID  int       `json:"teacher_user_id" example:"1"`
	TeacherEmail   string    `json:"teacher_email"   example:"test@test.com"`
	TeacherName    string    `json:"teacher_name"    example:"John"`
	TeacherSurname string    `json:"teacher_surname" example:"Smith"`
	TeacherAvatar  string    `json:"teacher_avatar"  example:"uuid.png"`
	CategoryID     int       `json:"category_id"     example:"1"`
	CategoryName   string    `json:"category_name"   example:"Programming"`
	Status         string    `json:"status"          example:"verification"`
	Datetime       time.Time `json:"datetime"        example:"2025-02-01T09:00:00Z"`
}
