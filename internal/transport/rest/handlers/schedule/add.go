package schedule

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
	"go.uber.org/zap"
)

const (
	AddRoute = "/schedule"
)

// AddScheduleTime returns http.HandlerFunc
// @Summary Add time to schedule
// @Description Add time to teacher schedule
// @Tags teachers
// @Accept json
// @Produce json
// @Param addTimeRequest body addTimeRequest true "datetime"
// @Success 201
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 403 {object} httputils.ErrorStruct
// @Failure 409 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /teacher/schedule [post]
// @Security     BearerAuth
func (h *ScheduleHandlers) AddScheduleTime() http.HandlerFunc {
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

		var req addTimeRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			if err = httputils.RespondWith400(w, "failed to decode body"); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		if req.Datetime.Before(time.Now()) {
			if err := httputils.RespondWith400(w, "the date must not be past"); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		err := h.scheduleService.AddTime(r.Context(), userID, req.Datetime)

		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				err = httputils.RespondWith401(w, err.Error())
			case errors.Is(err, serviceErrors.ErrorUserIsNotTeacher):
				err = httputils.RespondWithError(w, http.StatusForbidden, err.Error())
			case errors.Is(err, serviceErrors.ErrorScheduleTimeExists):
				err = httputils.RespondWithError(w, http.StatusConflict, err.Error())
			default:
				h.log.Error(err.Error())
				err = httputils.RespondWith500(w)
			}

			if err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}

			return
		}

		respondErr := httputils.SuccessRespondWith201(w, struct{}{})
		if respondErr != nil {
			h.log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

type addTimeRequest struct {
	Datetime time.Time `json:"datetime" example:"2025-02-01T00:00:00Z" binding:"required"`
}
