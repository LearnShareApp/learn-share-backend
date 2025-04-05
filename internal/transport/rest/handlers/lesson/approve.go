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
	approveRoute = "/{id}/approve"
)

// ApproveLesson returns http.HandlerFunc
// @Summary Approve lesson
// @Description Set lesson status "waiting" if this user is a teacher to lesson and lesson hasn't been cancelled (was verification)
// @Tags lessons
// @Produce json
// @Param id path int true "LessonID"
// @Success 200
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 403 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /lessons/{id}/approve [put]
// @Security     BearerAuth
func (h *LessonHandlers) ApproveLesson() http.HandlerFunc {
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

		err = h.lessonService.ApproveLesson(r.Context(), userID, lessonID)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				err = httputils.RespondWith401(w, err.Error())
			case errors.Is(err, serviceErrors.ErrorLessonNotFound):
				err = httputils.RespondWith404(w, err.Error())
			case errors.Is(err, serviceErrors.ErrorUserIsNotTeacher):
				err = httputils.RespondWith403(w, "unavailable operation for students")
			case errors.Is(err, serviceErrors.ErrorNotRelatedTeacherToLesson):
				err = httputils.RespondWith403(w, err.Error())
			case errors.Is(err, serviceErrors.ErrorStatusNonVerification):
				err = httputils.RespondWith403(w, "can approve a lesson if only the lesson had a verification status")
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
