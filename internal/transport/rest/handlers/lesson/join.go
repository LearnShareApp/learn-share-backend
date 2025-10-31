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
	joinRoute = "/{id}/join"
)

// JoinToLesson returns http.HandlerFunc
// @Summary Join the lesson
// @Description generate meet token to join "ongoing" lesson (if user related to lesson)
// @Tags lessons
// @Produce json
// @Param id path int true "LessonID"
// @Success 200 {object} connectLessonResponse
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 403 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /lessons/{id}/join [get]
// @Security     BearerAuth
func (h *LessonHandlers) JoinToLesson() http.HandlerFunc {
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

		token, err := h.lessonService.JoinLesson(r.Context(), userID, lessonID)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				httputils.RespondWith401(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorLessonNotFound):
				httputils.RespondWith404(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorNotRelatedUserToLesson):
				httputils.RespondWith403(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorUnavailableOperationState):
				httputils.RespondWith403(w, "can join to the lesson only if it is ongoing", h.log)
			default:
				h.log.Error(err.Error())
				httputils.RespondWith500(w, h.log)
			}

			return
		}

		resp := connectLessonResponse{Token: token}
		httputils.SuccessRespondWith200(w, resp, h.log)
	}
}

type connectLessonResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
