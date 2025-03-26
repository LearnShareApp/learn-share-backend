package become_teacher

import (
	"errors"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
	"go.uber.org/zap"
	"net/http"
)

const (
	Route = ""
)

// MakeHandler returns http.HandlerFunc
// @Summary User registrate also as teacher
// @Description Get user id by jwt token, and he became teach (if he was not be registrate himself as teacher)
// @Tags teachers
// @Produce json
// @Success 201
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 409 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /teacher [post]
// @Security     BearerAuth
func MakeHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.Context().Value(jwt.UserIDKey).(int)
		if id == 0 {
			log.Error("id was missed in context")
			if err := httputils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		err := s.Do(r.Context(), id)
		if err != nil {
			if errors.Is(err, serviceErrors.ErrorUserNotFound) {
				if err = httputils.RespondWith401(w, err.Error()); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return

			} else if errors.Is(err, serviceErrors.ErrorTeacherExists) {
				if err = httputils.RespondWithError(w, http.StatusConflict, err.Error()); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return

			} else {
				log.Error(err.Error())
				if err = httputils.RespondWith500(w); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return
			}
		}
		respondErr := httputils.SuccessRespondWith201(w, struct{}{})
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}
