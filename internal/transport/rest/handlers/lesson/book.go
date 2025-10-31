package lesson

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/httputils"
	"go.uber.org/zap"
)

const (
	bookRoute = "/lesson"
)

// BookLesson returns http.HandlerFunc
// @Summary Add new pending lesson (lesson request)
// @Description Check is all data confirmed and if so create lesson request (pending state)
// @Tags lessons
// @Accept json
// @Produce json
// @Param bookLessonRequest body bookLessonRequest true "LessonData"
// @Success 201
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 404 {object} httputils.ErrorStruct
// @Failure 409 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /lesson [post]
// @Security     BearerAuth
func (h *LessonHandlers) BookLesson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDValue := r.Context().Value(jwt.UserIDKey)
		userID, ok := userIDValue.(int)
		if !ok || userID == 0 {
			h.log.Error("invalid or missing user ID in context", zap.Any("value", userIDValue))
			httputils.RespondWith500(w, h.log)

			return
		}

		var req bookLessonRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httputils.RespondWith400(w, "failed to decode body", h.log)

			return
		}

		if req.TeacherID == 0 || req.CategoryID == 0 || req.ScheduleTimeID == 0 {
			httputils.RespondWith400(w, "teacher_id, category_id or schedule_time_id is empty (required)", h.log)

			return
		}

		lesson := entities.Lesson{
			StudentID:      userID,
			TeacherID:      req.TeacherID,
			CategoryID:     req.CategoryID,
			ScheduleTimeID: req.ScheduleTimeID,
		}

		err := h.lessonService.BookLesson(r.Context(), &lesson)

		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				httputils.RespondWith401(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorTeacherNotFound):
				httputils.RespondWith404(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorStudentAndTeacherSame):
				httputils.RespondWith400(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorCategoryNotFound):
				httputils.RespondWith404(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorSkillUnregistered):
				httputils.RespondWith404(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorSkillInactive):
				httputils.RespondWith404(w, "teacher's skill is inactive", h.log)
			case errors.Is(err, serviceErrors.ErrorScheduleTimeNotFound):
				httputils.RespondWith404(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorScheduleTimeForAnotherTeacher):
				httputils.RespondWith400(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorScheduleTimeUnavailable):
				httputils.RespondWith400(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorLessonTimeBooked):
				httputils.RespondWith409(w, err.Error(), h.log)
			default:
				h.log.Error(err.Error())
				httputils.RespondWith500(w, h.log)
			}

			return
		}

		httputils.SuccessRespondWith201(w, struct{}{}, h.log)
	}
}

// @Description book lesson body bookLessonRequest.
type bookLessonRequest struct {
	TeacherID      int `json:"teacher_id"       example:"1" binding:"required"` // @Description exactly teacherID, not his userID
	CategoryID     int `json:"category_id"      example:"1" binding:"required"`
	ScheduleTimeID int `json:"schedule_time_id" example:"1" binding:"required"`
}
