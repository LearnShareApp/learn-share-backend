package add_skill

import (
	"encoding/json"
	"errors"
	"net/http"

	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"github.com/LearnShareApp/learn-share-backend/internal/service/jwt"

	"go.uber.org/zap"
)

const (
	Route = "/skill"
)

// MakeHandler returns http.HandlerFunc
// @Summary Registrate new skill
// @Description Registrate new skill for teacher (if he not exists create and registrate skill)
// @Tags teachers
// @Accept json
// @Produce json
// @Param request body request true "Skill data"
// @Success 201
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 409 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /teacher/skill [post]
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

		if req.CategoryID == 0 {
			if err := httputils.RespondWith400(w, "category_id is empty (required)"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}

			return
		}

		err := s.Do(r.Context(), userId, req.CategoryID, req.VideoCardLink, req.About)
		if err != nil {
			if errors.Is(err, serviceErrors.ErrorUserNotFound) {
				if err = httputils.RespondWith401(w, serviceErrors.ErrorUserNotFound.Error()); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}

				return
			} else if errors.Is(err, serviceErrors.ErrorCategoryNotFound) {
				if err = httputils.RespondWith400(w, serviceErrors.ErrorCategoryNotFound.Error()); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}

				return
			} else if errors.Is(err, serviceErrors.ErrorSkillRegistered) {
				if err = httputils.RespondWithError(w, http.StatusConflict, serviceErrors.ErrorSkillRegistered.Error()); err != nil {
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
