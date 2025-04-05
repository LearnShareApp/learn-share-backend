package teacher

import (
	"encoding/json"
	"errors"
	"net/http"

	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
	"go.uber.org/zap"
)

const (
	addSkillRoute = "/skill"
)

// AddSkill returns http.HandlerFunc
// @Summary Registrate new skill
// @Description Registrate new skill for teacher (if he not exists create and registrate skill)
// @Tags teachers
// @Accept json
// @Produce json
// @Param addSkillRequest body addSkillRequest true "Skill data"
// @Success 201
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 409 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /teacher/skill [post]
// @Security     BearerAuth
func (h *TeacherHandlers) AddSkill() http.HandlerFunc {
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

		var req addSkillRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			if err = httputils.RespondWith400(w, "failed to decode body"); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}

			return
		}

		if req.CategoryID == 0 || req.About == "" || req.VideoCardLink == "" {
			var err error

			switch {
			case req.CategoryID == 0:
				err = httputils.RespondWith400(w, "category_id is empty (required)")
			case req.About == "":
				err = httputils.RespondWith400(w, "about is empty (required)")
			case req.VideoCardLink == "":
				err = httputils.RespondWith400(w, "video_card_link is empty (required)")
			}

			if err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}

			return
		}

		err := h.teacherService.AddSkill(r.Context(), userID, req.CategoryID, req.VideoCardLink, req.About)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				err = httputils.RespondWith401(w, err.Error())
			case errors.Is(err, serviceErrors.ErrorCategoryNotFound):
				err = httputils.RespondWith400(w, err.Error())
			case errors.Is(err, serviceErrors.ErrorSkillRegistered):
				err = httputils.RespondWithError(w, http.StatusConflict, err.Error())
			default:
				h.log.Error(err.Error())
				err = httputils.RespondWith500(w)
			}
			if err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}

			return
		}

		respondErr := httputils.SuccessRespondWith201(w, struct{}{})
		if respondErr != nil {
			h.log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

type addSkillRequest struct {
	CategoryID    int    `json:"category_id"     example:"1"                                                binding:"required"`
	VideoCardLink string `json:"video_card_link" example:"https://youtu.be/HIcSWuKMwOw?si=FtxN1QJU9ZWnXy85"`
	About         string `json:"about"           example:"I am Groot"`
}
