package user

import (
	"net/http"

	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/httputils"
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
			httputils.RespondWith500(w, h.log)

			return
		}

		isAdmin, err := h.userService.CheckUserOnAdminByID(r.Context(), userID)
		if err != nil {
			switch {
			default:
				h.log.Error("failed to check user on admin", zap.Error(err))
				httputils.RespondWith500(w, h.log)
			}

			return
		}

		resp := BoolResponse{IsAdmin: isAdmin}

		httputils.SuccessRespondWith200(w, resp, h.log)
	}
}

type BoolResponse struct {
	IsAdmin bool `json:"is_admin"`
}
