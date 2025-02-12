package get_times

import (
	"errors"
	"net/http"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"github.com/LearnShareApp/learn-share-backend/internal/service/jwt"
	"go.uber.org/zap"
)

const (
	PublicRoute    = "/{id}/schedule"
	ProtectedRoute = "/schedule"
)

// MakeProtectedHandler returns http.HandlerFunc
// @Summary Get times from schedule
// @Description Get lessons times from teacher schedule
// @Tags teachers
// @Produce json
// @Success 200 {object} response
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 404 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /teacher/schedule [get]
// @Security     BearerAuth
func MakeProtectedHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userId := r.Context().Value(jwt.UserIDKey).(int)
		if userId == 0 {
			log.Error("id was missed in context")
			if err := httputils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		times, err := s.DoByUserId(r.Context(), userId)

		if err != nil {
			coveringErrors(w, log, err)
			return
		}

		resp := mappingToResponse(times)

		respondErr := httputils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

// MakePublicHandler returns http.HandlerFunc which handle get schedule, get teacher id from http param
// @Summary Get times from schedule
// @Description Get lessons times from teacher schedule (by teacher ID)
// @Tags teachers
// @Produce json
// @Param id path int true "Teacher's ID"
// @Success 200 {object} response
// @Failure 403 {object} httputils.ErrorStruct
// @Failure 404 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /teachers/{id}/schedule [get]
func MakePublicHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teacherId, err := httputils.GetIntParamFromRequestPath(r, "id")

		if err != nil {
			if err := httputils.RespondWith400(w, "missed {id} param in url path"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		user, err := s.DoByTeacherId(r.Context(), teacherId)

		if err != nil {
			coveringErrors(w, log, err)
			return
		}

		resp := mappingToResponse(user)

		respondErr := httputils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

func coveringErrors(w http.ResponseWriter, log *zap.Logger, err error) {
	switch {
	case errors.Is(err, serviceErrors.ErrorUserIsNotTeacher):
		err = httputils.RespondWith403(w, serviceErrors.ErrorUserIsNotTeacher.Error())
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

func mappingToResponse(scheduleTimes []*entities.ScheduleTime) response {
	resp := response{
		Datetimes: make([]times, len(scheduleTimes)),
	}

	for i := range scheduleTimes {
		resp.Datetimes[i] = times{
			ScheduleTimeId: scheduleTimes[i].Id,
			Datetime:       scheduleTimes[i].Datetime,
			IsAvailable:    scheduleTimes[i].IsAvailable,
		}
	}

	return resp
}
