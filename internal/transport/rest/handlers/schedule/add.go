package schedule

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/httputils"
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
			httputils.RespondWith500(w, h.log)

			return
		}

		var req addTimeRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httputils.RespondWith400(w, "failed to decode body", h.log)

			return
		}

		if req.Datetime.Before(time.Now()) {
			httputils.RespondWith400(w, "the date must not be past", h.log)

			return
		}

		err := h.scheduleService.AddTime(r.Context(), userID, req.Datetime)

		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				httputils.RespondWith401(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorUserIsNotTeacher):
				httputils.RespondWith403(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorScheduleTimeExists):
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

type addTimeRequest struct {
	Datetime time.Time `json:"datetime" example:"2025-02-01T00:00:00Z" binding:"required"`
}
