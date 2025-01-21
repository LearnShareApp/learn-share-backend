package add_time

import (
	"encoding/json"
	"errors"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/jsonutils"
	"github.com/LearnShareApp/learn-share-backend/internal/service/jwt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	Route = "/schedule"
)

// MakeHandler returns http.HandlerFunc
// @Summary Add time to schedule
// @Description Add time to teacher schedule
// @Tags teachers
// @Accept json
// @Produce json
// @Param request body request true "datetime"
// @Success 201
// @Failure 400 {object} jsonutils.ErrorStruct
// @Failure 401 {object} jsonutils.ErrorStruct
// @Failure 403 {object} jsonutils.ErrorStruct
// @Failure 409 {object} jsonutils.ErrorStruct
// @Failure 500 {object} jsonutils.ErrorStruct
// @Router /teacher/schedule [post]
// @Security     BearerAuth
func MakeHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userId := r.Context().Value(jwt.UserIDKey).(int)
		if userId == 0 {
			log.Error("id was missed in context")
			if err := jsonutils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		var req request

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			if err = jsonutils.RespondWith400(w, "failed to decode body"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		if req.Datetime.Before(time.Now()) {
			if err := jsonutils.RespondWith400(w, "the date must not be past"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		err := s.Do(r.Context(), userId, req.Datetime)
		if err != nil {
			if errors.Is(err, serviceErrors.ErrorUserNotFound) {
				if err = jsonutils.RespondWith401(w, serviceErrors.ErrorUserNotFound.Error()); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return

			} else if errors.Is(err, serviceErrors.ErrorTeacherNotFound) {
				if err = jsonutils.RespondWithError(w, http.StatusForbidden, "you are not the teacher"); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return

			} else if errors.Is(err, serviceErrors.ErrorScheduleTimeExists) {
				if err = jsonutils.RespondWithError(w, http.StatusConflict, serviceErrors.ErrorScheduleTimeExists.Error()); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return

			} else {
				log.Error(err.Error())
				if err = jsonutils.RespondWith500(w); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return
			}
		}
		respondErr := jsonutils.SuccessRespondWith201(w, struct{}{})
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}
