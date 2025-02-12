package join_lesson

import (
	"errors"
	"net/http"
	"strconv"

	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"github.com/LearnShareApp/learn-share-backend/internal/service/jwt"
	"go.uber.org/zap"
)

const (
	Route = "/{id}/join"
)

// MakeHandler returns http.HandlerFunc
// @Summary Join the lesson
// @Description generate meet token to join "ongoing" lesson (if user related to lesson)
// @Tags lessons
// @Produce json
// @Param id path int true "LessonID"
// @Success 200 {object} response
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 403 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /lessons/{id}/join [get]
// @Security     BearerAuth
func MakeHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get lesson id from path
		paramId := r.PathValue("id")
		if paramId == "" {
			if err := httputils.RespondWith400(w, "missed {id} param in url path"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}
		lessonId, err := strconv.Atoi(paramId)
		if err != nil {
			log.Error("failed to parse id from URL path", zap.Error(err))
			if err := httputils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		// get userId from token
		userId := r.Context().Value(jwt.UserIDKey).(int)
		if userId == 0 {
			log.Error("id was missed in context")
			if err := httputils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		token, err := s.Do(r.Context(), userId, lessonId)
		if err != nil {
			if errors.Is(err, serviceErrors.ErrorUserNotFound) {
				if err = httputils.RespondWith401(w, serviceErrors.ErrorUserNotFound.Error()); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
			} else if errors.Is(err, serviceErrors.ErrorLessonNotFound) {
				if err = httputils.RespondWith404(w, serviceErrors.ErrorLessonNotFound.Error()); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
			} else if errors.Is(err, serviceErrors.ErrorUserIsNotTeacher) {
				if err = httputils.RespondWith403(w, "unavailable operation for students"); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
			} else if errors.Is(err, serviceErrors.ErrorNotRelatedUserToLesson) {
				if err = httputils.RespondWith403(w, serviceErrors.ErrorNotRelatedUserToLesson.Error()); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
			} else if errors.Is(err, serviceErrors.ErrorStatusNonOngoing) {
				if err = httputils.RespondWith403(w, "can start a lesson if only the lesson had a ongoing status"); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
			} else {
				log.Error(err.Error())
				if err = httputils.RespondWith500(w); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
			}
			return
		}

		resp := response{Token: token}
		respondErr := httputils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}
