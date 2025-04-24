package user

import (
	"net/http"

	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
	"go.uber.org/zap"
)

const checkOnAdminRoute = "/is-admin"

// CheckOnAdmin returns http.HandlerFunc
// @Summary Return boolean value is user an admin
// @Description Return boolean value is user an admin or not
// @Tags users
// @Produce json
// @Success 200 {object} BoolResponse
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /user/is-admin [get]
// @Security     BearerAuth
func (h *UserHandlers) CheckOnAdmin() http.HandlerFunc {
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

		isAdmin, err := h.userService.CheckUserOnAdminByID(r.Context(), userID)
		if err != nil {
			switch {
			default:
				h.log.Error("failed to check user on admin", zap.Error(err))
				err = httputils.RespondWith500(w)
			}

			if err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}

			return
		}

		resp := BoolResponse{IsAdmin: isAdmin}

		respondErr := httputils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			h.log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

type BoolResponse struct {
	IsAdmin bool `json:"is_admin"`
}
