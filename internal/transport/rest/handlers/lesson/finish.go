package lesson

import (
	"errors"
	"net/http"

	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/httputils"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
)

const (
	finishRoute = "/{id}/finish"
)

// FinishLesson returns http.HandlerFunc
// @Summary Finished lesson
// @Description Set lesson in finished state if this user is a teacher to lesson and lesson state has been ongoing
// @Tags lessons
// @Produce json
// @Param id path int true "LessonID"
// @Success 200
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 403 {object} httputils.ErrorStruct
// @Failure 404 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /lessons/{id}/finish [put]
// @Security     BearerAuth
func (h *LessonHandlers) FinishLesson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get lesson id from path
		lessonID, err := httputils.GetIntParamFromRequestPath(r, "id")
		if err != nil {
			httputils.RespondWith400(w, "missed {id} param in url path", h.log)

			return
		}

		// get userID from token
		userID := r.Context().Value(jwt.UserIDKey).(int)
		if userID == 0 {
			h.log.Error("id was missed in context")
			httputils.RespondWith500(w, h.log)

			return
		}

		err = h.lessonService.FinishLesson(r.Context(), userID, lessonID)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				httputils.RespondWith401(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorLessonNotFound):
				httputils.RespondWith404(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorUserIsNotTeacher):
				httputils.RespondWith403(w, "unavailable operation for students", h.log)
			case errors.Is(err, serviceErrors.ErrorNotRelatedTeacherToLesson):
				httputils.RespondWith403(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorUnavailableStateTransition):
				httputils.RespondWith403(w, "can finish a lesson if only the lesson had been ongoing", h.log)
			default:
				h.log.Error(err.Error())
				httputils.RespondWith500(w, h.log)
			}

			return
		}

		httputils.SuccessRespondWith200(w, struct{}{}, h.log)
	}
}
