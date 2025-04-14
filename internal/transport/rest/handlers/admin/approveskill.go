package admin

import (
	"errors"
	"net/http"

	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
	"go.uber.org/zap"
)

const approveSkillRoute = "/skills/{id}/approve"

// ApproveSkill returns http.HandlerFunc
// @Summary approve teacher'skill
// @Description handler for approving teacher's skill
// @Tags admin
// @Produce json
// @Param id path int true "skillID"
// @Success 200
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 403 {object} httputils.ErrorStruct
// @Failure 404 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /admin/skills/{id}/approve [put]
// @Security     BearerAuth
func (h *AdminHandlers) ApproveSkill() http.HandlerFunc {
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

		// get skill id from path
		skillID, err := httputils.GetIntParamFromRequestPath(r, "id")

		if err != nil {
			if err := httputils.RespondWith400(w, "missed {id} param in url path"); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		isAdmin, err := h.service.CheckUserOnAdminByID(r.Context(), userID)

		if err != nil {
			h.log.Error("failed to check user on admin", zap.Error(err))
			if err := httputils.RespondWith500(w); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}

			return
		}

		if !isAdmin {
			if err := httputils.RespondWith403(w, serviceErrors.ErrorNotAdmin.Error()); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}

			return
		}

		err = h.service.ApproveTeacherSkill(r.Context(), skillID)

		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorSkillNotFound):
				err = httputils.RespondWith404(w, err.Error())
			case errors.Is(err, serviceErrors.ErrorSkillAlreadyApproved):
				err = httputils.RespondWith409(w, err.Error())
			default:
				h.log.Error(err.Error())
				err = httputils.RespondWith500(w)
			}

			if err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}

			return
		}

		if err = httputils.SuccessRespondWith200(w, struct{}{}); err != nil {
			h.log.Error("failed to send response", zap.Error(err))
		}
	}
}
