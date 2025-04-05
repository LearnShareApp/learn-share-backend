package lesson

import (
	"errors"
	"net/http"

	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"go.uber.org/zap"
)

const (
	getShortRoute = "/{id}/short-data"
)

// GetLessonShortData returns http.HandlerFunc
// @Summary Get lesson really short data by lesson's id
// @Description Return lesson short data by lesson's id
// @Tags lessons
// @Produce json
// @Param id path int true "LessonID"
// @Success 200 {object} getLessonShortDataResponse
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 404 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /lessons/{id}/short-data [get]
func (h *LessonHandlers) GetLessonShortData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get lesson id from path
		lessonID, err := httputils.GetIntParamFromRequestPath(r, "id")
		if err != nil {
			if err := httputils.RespondWith400(w, "missed {id} param in url path"); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		lesson, err := h.lessonService.GetLessonShortData(r.Context(), lessonID)

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

		resp := getLessonShortDataResponse{
			LessonID: lessonID,

			TeacherID:     lesson.TeacherID,
			TeacherUserID: lesson.TeacherUserData.Id,

			StudentID: lesson.StudentID,

			CategoryID:   lesson.CategoryID,
			CategoryName: lesson.CategoryName,
		}

		respondErr := httputils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			h.log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

type getLessonShortDataResponse struct {
	LessonID      int    `json:"lesson_id"       example:"1"`
	TeacherID     int    `json:"teacher_id"      example:"1"`
	TeacherUserID int    `json:"teacher_user_id" example:"1"`
	StudentID     int    `json:"student_id"      example:"1"`
	CategoryID    int    `json:"category_id"     example:"1"`
	CategoryName  string `json:"category_name"   example:"Programming"`
}
