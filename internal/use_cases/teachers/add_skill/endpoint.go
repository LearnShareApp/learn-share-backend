package add_skill

import (
	"encoding/json"
	"errors"
	"net/http"

	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"

	"github.com/LearnShareApp/learn-share-backend/internal/jsonutils"
	"github.com/LearnShareApp/learn-share-backend/internal/service/jwt"
	"go.uber.org/zap"
)

const (
	Route = "/skill"
)

// MakeHandler returns http.HandlerFunc
// @Summary Registrate new skill
// @Description Registrate new skill for teacher (if he not exists create and registrate skill)
// @Tags teacher
// @Accept json
// @Produce json
// @Param request body request true "Skill data"
// @Success 201
// @Failure 400 {object} jsonutils.ErrorStruct
// @Failure 401 {object} jsonutils.ErrorStruct
// @Failure 409 {object} jsonutils.ErrorStruct
// @Failure 500 {object} jsonutils.ErrorStruct
// @Router /teacher/skill [post]
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

		if req.CategoryId == 0 {
			if err := jsonutils.RespondWith400(w, "category_id is empty (required)"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		err := s.Do(r.Context(), userId, req.CategoryId, req.VideoCardLink, req.About)
		if err != nil {
			if errors.Is(err, serviceErrors.ErrorUserNotFound) {
				if err = jsonutils.RespondWith401(w, serviceErrors.ErrorUserNotFound.Error()); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return

			} else if errors.Is(err, serviceErrors.ErrorCategoryNotFound) {
				if err = jsonutils.RespondWith400(w, serviceErrors.ErrorCategoryNotFound.Error()); err != nil {
					log.Error("failed to send response", zap.Error(err))
				}
				return
			} else if errors.Is(err, serviceErrors.ErrorSkillRegistered) {
				if err = jsonutils.RespondWithError(w, http.StatusConflict, serviceErrors.ErrorSkillRegistered.Error()); err != nil {
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
