package teacher

import (
	"errors"
	"net/http"

	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
	"go.uber.org/zap"
)

const (
	becomeRoute = "/"
)

// BecomeTeacher returns http.HandlerFunc
// @Summary User registrate also as teacher
// @Description Get user id by jwt token, and he became teach (if he was not be registrate himself as teacher)
// @Tags teachers
// @Produce json
// @Success 201
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 409 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /teacher [post]
// @Security     BearerAuth
func (h *TeacherHandlers) BecomeTeacher() http.HandlerFunc {
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

		err := h.teacherService.BecomeTeacher(r.Context(), userID)

		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				err = httputils.RespondWith401(w, err.Error())
			case errors.Is(err, serviceErrors.ErrorTeacherExists):
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
