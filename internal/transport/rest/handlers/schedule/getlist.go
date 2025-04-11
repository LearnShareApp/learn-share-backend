package schedule

import (
	"errors"
	"net/http"
	"time"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
	"go.uber.org/zap"
)

const (
	PublicGetListRoute    = "/{id}/schedule"
	ProtectedGetListRoute = "/schedule"
)

// GetScheduleProtected returns http.HandlerFunc
// @Summary Get times from schedule
// @Description Get lessons times from teacher schedule
// @Tags teachers
// @Produce json
// @Success 200 {object} getTimesResponse
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 404 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /teacher/schedule [get]
// @Security     BearerAuth
func (h *ScheduleHandlers) GetScheduleProtected() http.HandlerFunc {
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

		teacher, err := h.scheduleService.GetTeacherByUserID(r.Context(), userID)

		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorTeacherNotFound):
				err = httputils.RespondWith403(w, serviceErrors.ErrorUserIsNotTeacher.Error())
			default:
				h.log.Error(err.Error())
				err = httputils.RespondWith500(w)
			}

			if err != nil {
				h.log.Error("failed to send error response", zap.Error(err))
			}

			return
		}

		times, err := h.scheduleService.GetTimes(r.Context(), teacher)

		if err != nil {
			coveringErrors(w, h.log, err)
			return
		}

		resp := mappingToResponse(times)

		respondErr := httputils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			h.log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

// GetSchedulePublic returns http.HandlerFunc which handle get schedule, get teacher id from http param
// @Summary Get times from schedule
// @Description Get lessons times from teacher schedule (by teacher ID)
// @Tags teachers
// @Produce json
// @Param id path int true "Teacher's ID"
// @Success 200 {object} getTimesResponse
// @Failure 403 {object} httputils.ErrorStruct
// @Failure 404 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /teachers/{id}/schedule [get]
func (h *ScheduleHandlers) GetSchedulePublic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teacherID, err := httputils.GetIntParamFromRequestPath(r, "id")

		if err != nil {
			if err := httputils.RespondWith400(w, "missed {id} param in url path"); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		teacher, err := h.scheduleService.GetTeacherByID(r.Context(), teacherID)

		if err != nil {
			coveringErrors(w, h.log, err)
			return
		}

		times, err := h.scheduleService.GetTimes(r.Context(), teacher)

		if err != nil {
			coveringErrors(w, h.log, err)
			return
		}

		resp := mappingToResponse(times)

		respondErr := httputils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			h.log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

func coveringErrors(w http.ResponseWriter, log *zap.Logger, err error) {
	switch {
	case errors.Is(err, serviceErrors.ErrorTeacherNotFound):
		err = httputils.RespondWith404(w, serviceErrors.ErrorTeacherNotFound.Error())
	default:
		log.Error(err.Error())
		err = httputils.RespondWith500(w)
	}

	if err != nil {
		log.Error("failed to send error response", zap.Error(err))
	}
}

func mappingToResponse(scheduleTimes []*entities.ScheduleTime) getTimesResponse {
	resp := getTimesResponse{
		Datetimes: make([]respTimes, len(scheduleTimes)),
	}

	for i := range scheduleTimes {
		resp.Datetimes[i] = respTimes{
			ScheduleTimeID: scheduleTimes[i].ID,
			Datetime:       scheduleTimes[i].Datetime,
			IsAvailable:    scheduleTimes[i].IsAvailable,
		}
	}

	return resp
}

type getTimesResponse struct {
	Datetimes []respTimes `json:"datetimes"`
}

type respTimes struct {
	ScheduleTimeID int       `json:"schedule_time_id" example:"1"`
	Datetime       time.Time `json:"datetime"         example:"0001-01-01T00:00:00Z"`
	IsAvailable    bool      `json:"is_available"     example:"true"`
}
