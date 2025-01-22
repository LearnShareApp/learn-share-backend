package get_times

import (
	"errors"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/jsonutils"
	"github.com/LearnShareApp/learn-share-backend/internal/service/jwt"
	"go.uber.org/zap"
	"net/http"
	"strconv"
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
// @Failure 401 {object} jsonutils.ErrorStruct
// @Failure 404 {object} jsonutils.ErrorStruct
// @Failure 500 {object} jsonutils.ErrorStruct
// @Router /teacher/schedule [get]
// @Security     BearerAuth
func MakeProtectedHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userId := r.Context().Value(jwt.UserIDKey).(int)
		if userId == 0 {
			log.Error("id was missed in context")
			if err := jsonutils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		times, err := s.Do(r.Context(), userId)

		if err != nil {
			coveringErrors(w, log, err)
			return
		}

		resp := mappingToResponse(times)

		respondErr := jsonutils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

// MakePublicHandler returns http.HandlerFunc which handle get schedule, get user id from http param
// @Summary Get times from schedule
// @Description Get lessons times from teacher schedule (by his UserID)
// @Tags teachers
// @Produce json
// @Param id path int true "Teacher's UserID"
// @Success 200 {object} response
// @Failure 404 {object} jsonutils.ErrorStruct
// @Failure 500 {object} jsonutils.ErrorStruct
// @Router /teachers/{id}/schedule [get]
func MakePublicHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id int

		paramId := r.PathValue("id")
		if paramId == "" {
			if err := jsonutils.RespondWith400(w, "missed {id} param in url path"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		id, err := strconv.Atoi(paramId)

		if err != nil {
			log.Error("failed to parse id from URL path", zap.Error(err))
			if err := jsonutils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		user, err := s.Do(r.Context(), id)

		if err != nil {
			coveringErrors(w, log, err)
			return
		}

		resp := mappingToResponse(user)

		respondErr := jsonutils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

func coveringErrors(w http.ResponseWriter, log *zap.Logger, err error) {
	if errors.Is(err, serviceErrors.ErrorUserNotFound) {
		if err := jsonutils.RespondWith404(w, serviceErrors.ErrorUserNotFound.Error()); err != nil {
			log.Error("failed to send response", zap.Error(err))
		}
	} else if errors.Is(err, serviceErrors.ErrorUserIsNotTeacher) {
		if err := jsonutils.RespondWith404(w, serviceErrors.ErrorUserIsNotTeacher.Error()); err != nil {
			log.Error("failed to send response", zap.Error(err))
		}
	} else {
		log.Error(err.Error())
		if err = jsonutils.RespondWith500(w); err != nil {
			log.Error("failed to send response", zap.Error(err))
		}
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
