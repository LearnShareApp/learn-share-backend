package lesson

import (
	"errors"
	"net/http"
	"time"

	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/httputils"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
	"go.uber.org/zap"
)

const (
	getForTeacherListRoute = "/teacher/lessons"
)

// GetForTeacherList returns http.HandlerFunc
// @Summary Get lessons for teachers
// @Description Return all lessons which have teacher
// @Tags teachers
// @Produce json
// @Success 200 {object} getTeacherLessonsResponse
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 403 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /teacher/lessons [get]
// @Security     BearerAuth
func (h *LessonHandlers) GetForTeacherList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDValue := r.Context().Value(jwt.UserIDKey)
		userID, ok := userIDValue.(int)
		if !ok || userID == 0 {
			h.log.Error("invalid or missing user ID in context", zap.Any("value", userIDValue))
			httputils.RespondWith500(w, h.log)

			return
		}

		lessons, err := h.lessonService.GetTeacherLessonList(r.Context(), userID)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				httputils.RespondWith401(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorUserIsNotTeacher):
				httputils.RespondWith403(w, err.Error(), h.log)
			default:
				h.log.Error(err.Error())
				httputils.RespondWith500(w, h.log)
			}

			return
		}

		var resp getTeacherLessonsResponse

		if lessons != nil {
			resp = getTeacherLessonsResponse{
				Lessons: make([]respTeacherLessons, len(lessons)),
			}

			for i := range lessons {
				resp.Lessons[i] = respTeacherLessons{
					LessonID:       lessons[i].ID,
					StudentID:      lessons[i].StudentID,
					StudentEmail:   lessons[i].StudentUserData.Email,
					StudentName:    lessons[i].StudentUserData.Name,
					StudentSurname: lessons[i].StudentUserData.Surname,
					StudentAvatar:  lessons[i].StudentUserData.Avatar,
					CategoryID:     lessons[i].CategoryID,
					CategoryName:   lessons[i].CategoryName,
					StateID:        lessons[i].StateMachineItem.StateID,
					StateName:      lessons[i].StateMachineItem.StateName,
					Status:         lessons[i].StatusName,
					Datetime:       lessons[i].ScheduleTimeDatetime,
				}
			}
		}

		httputils.SuccessRespondWith200(w, resp, h.log)
	}
}

type getTeacherLessonsResponse struct {
	Lessons []respTeacherLessons `json:"lessons"`
}

type respTeacherLessons struct {
	LessonID       int       `json:"lesson_id"       example:"1"`
	StudentID      int       `json:"student_id"      example:"1"`
	StudentEmail   string    `json:"student_email"   example:"test@test.com"`
	StudentName    string    `json:"student_name"    example:"John"`
	StudentSurname string    `json:"student_surname" example:"Smith"`
	StudentAvatar  string    `json:"student_avatar"  example:"uuid.png"`
	CategoryID     int       `json:"category_id"     example:"1"`
	CategoryName   string    `json:"category_name"   example:"Programming"`
	StateID        int       `json:"state_id"        example:"1"`
	StateName      string    `json:"state_name"      example:"pending"`
	Status         string    `json:"status"          example:"verification"`
	Datetime       time.Time `json:"datetime"        example:"2025-02-01T09:00:00Z"`
}
