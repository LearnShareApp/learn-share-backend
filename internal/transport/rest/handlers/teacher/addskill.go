package teacher

import (
	"encoding/json"
	"errors"
	"net/http"

	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/httputils"
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
			httputils.RespondWith500(w, h.log)

			return
		}

		var req addSkillRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httputils.RespondWith400(w, "failed to decode body", h.log)

			return
		}

		if req.CategoryID == 0 || req.About == "" || req.VideoCardLink == "" {
			switch {
			case req.CategoryID == 0:
				httputils.RespondWith400(w, "category_id is empty (required)", h.log)
			case req.About == "":
				httputils.RespondWith400(w, "about is empty (required)", h.log)
			case req.VideoCardLink == "":
				httputils.RespondWith400(w, "video_card_link is empty (required)", h.log)
			}

			return
		}

		err := h.teacherService.AddSkill(r.Context(), userID, req.CategoryID, req.VideoCardLink, req.About)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				httputils.RespondWith401(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorCategoryNotFound):
				httputils.RespondWith400(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorSkillRegistered):
				httputils.RespondWith409(w, err.Error(), h.log)
			default:
				h.log.Error(err.Error())
				httputils.RespondWith500(w, h.log)
			}

			return
		}

		httputils.SuccessRespondWith201(w, struct{}{}, h.log)
	}
}

type addSkillRequest struct {
	CategoryID    int    `json:"category_id"     example:"1"                                                binding:"required"`
	VideoCardLink string `json:"video_card_link" example:"https://youtu.be/HIcSWuKMwOw?si=FtxN1QJU9ZWnXy85"`
	About         string `json:"about"           example:"I am Groot"`
}
