package book_lesson

import (
	"encoding/json"
	"errors"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
	"net/http"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"go.uber.org/zap"
)

const (
	Route = ""
)

// MakeHandler returns http.HandlerFunc
// @Summary Add Unconfirmed lesson (lesson request)
// @Description Check is all data confirmed and if so create lesson with status "verification" (Unconfirmed)
// @Tags lessons
// @Accept json
// @Produce json
// @Param request body request true "LessonData"
// @Success 201
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 409 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /lesson [post]
// @Security     BearerAuth
func MakeHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userId := r.Context().Value(jwt.UserIDKey).(int)
		if userId == 0 {
			log.Error("id was missed in context")
			if err := httputils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		var req request

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			if err = httputils.RespondWith400(w, "failed to decode body"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		if req.TeacherId == 0 || req.CategoryId == 0 || req.ScheduleTimeId == 0 {
			if err := httputils.RespondWith400(w, "teacher_id, category_id or schedule_time_id is empty (required)"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		lesson := entities.Lesson{
			StudentId:      userId,
			TeacherId:      req.TeacherId,
			CategoryId:     req.CategoryId,
			ScheduleTimeId: req.ScheduleTimeId,
		}

		err := s.Do(r.Context(), &lesson)
		if err != nil {
			if errors.Is(err, serviceErrors.ErrorUserNotFound) {
				if err = httputils.RespondWith401(w, serviceErrors.ErrorUserNotFound.Error()); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return

			} else if errors.Is(err, serviceErrors.ErrorTeacherNotFound) {
				if err = httputils.RespondWith400(w, "teacher with such teacher_id not found"); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return

			} else if errors.Is(err, serviceErrors.ErrorStudentAndTeacherSame) {
				if err = httputils.RespondWith400(w, serviceErrors.ErrorStudentAndTeacherSame.Error()); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return

			} else if errors.Is(err, serviceErrors.ErrorCategoryNotFound) {
				if err = httputils.RespondWith400(w, "category with such category id not found"); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return

			} else if errors.Is(err, serviceErrors.ErrorSkillUnregistered) {
				if err = httputils.RespondWith400(w, "this teacher has not such skill"); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return

			} else if errors.Is(err, serviceErrors.ErrorScheduleTimeNotFound) {
				if err = httputils.RespondWith400(w, "schedule time with such schedule_time_id not found"); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return

			} else if errors.Is(err, serviceErrors.ErrorScheduleTimeForAnotherTeacher) {
				if err = httputils.RespondWith400(w, "schedule time belongs to another teacher"); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return

			} else if errors.Is(err, serviceErrors.ErrorScheduleTimeUnavailable) {
				if err = httputils.RespondWith400(w, serviceErrors.ErrorScheduleTimeUnavailable.Error()); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return

			} else if errors.Is(err, serviceErrors.ErrorLessonTimeBooked) {
				if err = httputils.RespondWithError(w, http.StatusConflict, serviceErrors.ErrorLessonTimeBooked.Error()); err != nil {
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
