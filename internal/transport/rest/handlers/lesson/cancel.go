package lesson

import (
	"errors"
	"net/http"

	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/httputils"
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
			httputils.RespondWith500(w, h.log)

			return
		}

		// get lesson id from path
		lessonID, err := httputils.GetIntParamFromRequestPath(r, "id")
		if err != nil {
			httputils.RespondWith400(w, "missed {id} param in url path", h.log)

			return
		}

		err = h.lessonService.CancelLesson(r.Context(), userID, lessonID)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				httputils.RespondWith401(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorLessonNotFound):
				httputils.RespondWith404(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorNotRelatedUserToLesson):
				httputils.RespondWith403(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorUserIsNotTeacher):
				httputils.RespondWith403(w, "cancel ongoing lesson is unavailable for student", h.log)
			case errors.Is(err, serviceErrors.ErrorNotRelatedTeacherToLesson):
				httputils.RespondWith403(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorUnavailableStateTransition):
				httputils.RespondWith403(w, "cannot cancel lesson, unavailable state transition", h.log)
			default:
				h.log.Error(err.Error())
				httputils.RespondWith500(w, h.log)
			}

			return
		}

		httputils.SuccessRespondWith200(w, struct{}{}, h.log)

	}
}
