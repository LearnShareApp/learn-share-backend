package admin

import (
	"errors"
	"net/http"

	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/httputils"
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
			httputils.RespondWith500(w, h.log)
			return
		}

		// get skill id from path
		skillID, err := httputils.GetIntParamFromRequestPath(r, "id")
		if err != nil {
			httputils.RespondWith400(w, "missed {id} param in url path", h.log)
			return
		}

		isAdmin, err := h.service.CheckUserOnAdminByID(r.Context(), userID)
		if err != nil {
			h.log.Error("failed to check user on admin", zap.Error(err))
			httputils.RespondWith500(w, h.log)
			return
		}

		if !isAdmin {
			httputils.RespondWith403(w, serviceErrors.ErrorNotAdmin.Error(), h.log)
			return
		}

		err = h.service.ApproveTeacherSkill(r.Context(), skillID)

		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorSkillNotFound):
				httputils.RespondWith404(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorSkillAlreadyApproved):
				httputils.RespondWith409(w, err.Error(), h.log)
			default:
				h.log.Error(err.Error())
				httputils.RespondWith500(w, h.log)
			}
			return
		}

		httputils.SuccessRespondWith200(w, struct{}{}, h.log)
	}
}
