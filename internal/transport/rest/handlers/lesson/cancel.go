package lesson

import (
	"errors"
	"net/http"

	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
	"go.uber.org/zap"
)

const (
	cancelRoute = "/{id}/cancel"
)

// CancelLesson returns http.HandlerFunc
// @Summary Cancel lesson
// @Description Set lesson in cancelled state if this user related to lesson
// @Tags lessons
// @Produce json
// @Param id path int true "LessonID"
// @Success 200
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 403 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /lessons/{id}/cancel [put]
// @Security     BearerAuth
func (h *LessonHandlers) CancelLesson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get userID from token
		userIDValue := r.Context().Value(jwt.UserIDKey)
		userID, ok := userIDValue.(int)
		if !ok || userID == 0 {
			h.log.Error("invalid or missing user ID in context", zap.Any("value", userIDValue))
			if err := httputils.RespondWith500(w); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		// get lesson id from path
		lessonID, err := httputils.GetIntParamFromRequestPath(r, "id")
		if err != nil {
			if err := httputils.RespondWith400(w, "missed {id} param in url path"); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		err = h.lessonService.CancelLesson(r.Context(), userID, lessonID)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				err = httputils.RespondWith401(w, err.Error())
			case errors.Is(err, serviceErrors.ErrorLessonNotFound):
				err = httputils.RespondWith404(w, err.Error())
			case errors.Is(err, serviceErrors.ErrorNotRelatedUserToLesson):
				err = httputils.RespondWith403(w, err.Error())
			case errors.Is(err, serviceErrors.ErrorUserIsNotTeacher):
				err = httputils.RespondWith403(w, "cancel ongoing lesson is unavailable for student")
			case errors.Is(err, serviceErrors.ErrorNotRelatedTeacherToLesson):
				err = httputils.RespondWith403(w, err.Error())
			case errors.Is(err, serviceErrors.ErrorUnavailableStateTransition):
				err = httputils.RespondWith403(w, "cannot cancel lesson, unavailable state transition")
			default:
				h.log.Error(err.Error())
				err = httputils.RespondWith500(w)
			}

			if err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		respondErr := httputils.SuccessRespondWith200(w, struct{}{})
		if respondErr != nil {
			h.log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}
