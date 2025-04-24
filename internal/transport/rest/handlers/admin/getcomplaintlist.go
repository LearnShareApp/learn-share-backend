package admin

import (
	"net/http"
	"time"

	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
	"go.uber.org/zap"
)

const getComplaintListRoute = "/complaints"

// GetAllComplaintList returns http.HandlerFunc
// @Summary get complaint's list
// @Description returns the list of complaints
// @Tags admin
// @Produce json
// @Success 200 {object} getComplaintListResponse
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 403 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /admin/complaints [get]
// @Security     BearerAuth
func (h *AdminHandlers) GetAllComplaintList() http.HandlerFunc {
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

		complaints, err := h.service.GetComplaintList(r.Context())

		if err != nil {
			switch {
			default:
				h.log.Error(err.Error())
				err = httputils.RespondWith500(w)
			}

			if err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}

			return
		}

		resp := getComplaintListResponse{
			Complaints: make([]respComplaint, 0, len(complaints)),
		}

		for i := range complaints {
			resp.Complaints = append(resp.Complaints, respComplaint{
				ComplaintID:       complaints[i].ComplaintID,
				ComplainerID:      complaints[i].ComplainerID,
				ComplainerName:    complaints[i].Complainer.Name,
				ComplainerSurname: complaints[i].Complainer.Surname,
				ComplainerEmail:   complaints[i].Complainer.Email,
				ComplainerAvatar:  complaints[i].Complainer.Avatar,
				ReportedID:        complaints[i].ReportedID,
				ReportedName:      complaints[i].Reported.Name,
				ReportedSurname:   complaints[i].Reported.Surname,
				ReportedEmail:     complaints[i].Reported.Email,
				ReportedAvatar:    complaints[i].Reported.Avatar,
				Reason:            complaints[i].Reason,
				Description:       complaints[i].Description,
				Date:              complaints[i].CreatedAt,
			})
		}

		if err = httputils.SuccessRespondWith200(w, resp); err != nil {
			h.log.Error("failed to send response", zap.Error(err))
		}
	}
}

type getComplaintListResponse struct {
	Complaints []respComplaint `json:"complaints"`
}

type respComplaint struct {
	ComplaintID int `json:"complaint_id"       example:"1"`

	ComplainerID      int    `json:"complainer_id"      example:"1"`
	ComplainerName    string `json:"complainer_name"    example:"John"`
	ComplainerSurname string `json:"complainer_surname" example:"Smith"`
	ComplainerEmail   string `json:"complainer_email"   example:"test@test.com"`
	ComplainerAvatar  string `json:"complainer_avatar"  example:"uuid.png"`

	ReportedID      int    `json:"reported_id"      example:"1"`
	ReportedName    string `json:"reported_name"    example:"John"`
	ReportedSurname string `json:"reported_surname" example:"Smith"`
	ReportedEmail   string `json:"reported_email"   example:"test@test.com"`
	ReportedAvatar  string `json:"reported_avatar"  example:"uuid.png"`

	Reason      string    `json:"reason"      example:"reason"`
	Description string    `json:"description" example:"description"`
	Date        time.Time `json:"date"        example:"2025-01-09T10:10:10+09:00"`
}
