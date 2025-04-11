package lesson

import (
	"errors"
	"net/http"
	"time"

	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"go.uber.org/zap"
)

const (
	getRoute = "/{id}"
)

// GetLesson returns http.HandlerFunc
// @Summary Get lesson data by lesson's id
// @Description Return lesson data by lesson's id
// @Tags lessons
// @Produce json
// @Param id path int true "LessonID"
// @Success 200 {object} getLessonResponse
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 404 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /lessons/{id} [get]
func (h *LessonHandlers) GetLesson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get lesson id from path
		lessonID, err := httputils.GetIntParamFromRequestPath(r, "id")
		if err != nil {
			if err := httputils.RespondWith400(w, "missed {id} param in url path"); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		lesson, err := h.lessonService.GetLesson(r.Context(), lessonID)

		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorLessonNotFound):
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

		resp := getLessonResponse{
			LessonID: lessonID,

			TeacherID:      lesson.TeacherID,
			TeacherUserID:  lesson.TeacherUserData.ID,
			TeacherEmail:   lesson.TeacherUserData.Email,
			TeacherName:    lesson.TeacherUserData.Name,
			TeacherSurname: lesson.TeacherUserData.Surname,
			TeacherAvatar:  lesson.TeacherUserData.Avatar,

			StudentID:      lesson.StudentID,
			StudentEmail:   lesson.StudentUserData.Email,
			StudentName:    lesson.StudentUserData.Name,
			StudentSurname: lesson.StudentUserData.Surname,
			StudentAvatar:  lesson.StudentUserData.Avatar,

			CategoryID:   lesson.CategoryID,
			CategoryName: lesson.CategoryName,
			Status:       lesson.StatusName,
			Datetime:     lesson.ScheduleTimeDatetime,
		}

		respondErr := httputils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			h.log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

// @Description data about lesson getLessonResponse.
type getLessonResponse struct {
	LessonID int `json:"lesson_id" example:"1"`

	TeacherID      int    `json:"teacher_id"      example:"1"`
	TeacherUserID  int    `json:"teacher_user_id" example:"1"`
	TeacherEmail   string `json:"teacher_email"   example:"test@test.com"`
	TeacherName    string `json:"teacher_name"    example:"John"`
	TeacherSurname string `json:"teacher_surname" example:"Smith"`
	TeacherAvatar  string `json:"teacher_avatar"  example:"uuid.png"`

	StudentID      int    `json:"student_id"      example:"1"`
	StudentEmail   string `json:"student_email"   example:"test@test.com"`
	StudentName    string `json:"student_name"    example:"John"`
	StudentSurname string `json:"student_surname" example:"Smith"`
	StudentAvatar  string `json:"student_avatar"  example:"uuid.png"`

	CategoryID   int       `json:"category_id"   example:"1"`
	CategoryName string    `json:"category_name" example:"Programming"`
	Status       string    `json:"status"        example:"verification"`
	Datetime     time.Time `json:"datetime"      example:"2025-02-01T09:00:00Z"`
}
