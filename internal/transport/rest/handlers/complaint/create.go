package complaint

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
)

const (
	createRoute = "/"
)

// CreateComplaint returns http.HandlerFunc
// @Summary Create a new complaint
// @Description Creating a new complaint to user (reported_id => user_id which you would like to report)
// @Tags complaint
// @Accept json
// @Produce json
// @Param createComplaintRequest body createComplaintRequest true "ComplaintData"
// @Success 201
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 404 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /complaint [post]
// @Security     BearerAuth
func (h *ComplaintHandlers) CreateComplaint() http.HandlerFunc {
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

		var req createComplaintRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			if err = httputils.RespondWith400(w, "failed to decode body"); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		if req.ReportedID == 0 || req.Reason == "" || req.Description == "" {
			if err := httputils.RespondWith400(w, "reported_id, reason or description is empty (required)"); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		complaint := entities.Complaint{
			ComplainerID: userID,
			ReportedID:   req.ReportedID,
			Reason:       req.Reason,
			Description:  req.Description,
		}

		err := h.service.CreateComplaint(r.Context(), &complaint)

		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				err = httputils.RespondWith401(w, err.Error())
			case errors.Is(err, serviceErrors.ErrorReportedUserNotFound):
				err = httputils.RespondWith404(w, err.Error())
			case errors.Is(err, serviceErrors.ErrorComplainerAndReportedSame):
				err = httputils.RespondWith400(w, err.Error())
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

type createComplaintRequest struct {
	ReportedID  int    `json:"reported_id" example:"1"                   binding:"required"`
	Reason      string `json:"reason"      example:"Rude attitude"       binding:"required"`
	Description string `json:"description" example:"your description..." binding:"required"`
}
