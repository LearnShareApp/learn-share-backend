package review

import (
	"errors"
	"net/http"

	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"go.uber.org/zap"
)

const (
	getListRoute = "/teachers/{id}/reviews"
)

// GetReviewList returns http.HandlerFunc
// @Summary Get reviews
// @Description Get all reviews about teacher
// @Tags teachers
// @Produce json
// @Param id path int true "Teacher ID"
// @Success 200 {object} getReviewResponse
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 404 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /teachers/{id}/reviews [get]
func (h *ReviewHandlers) GetReviewList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teacherID, err := httputils.GetIntParamFromRequestPath(r, "id")
		if err != nil {
			h.log.Error("failed to parse id from URL path", zap.Error(err))
			if err := httputils.RespondWith400(w, "missed {id} param in url path"); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		reviews, err := h.reviewService.GetReviews(r.Context(), teacherID)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorTeacherNotFound):
				err = httputils.RespondWith404(w, err.Error())
			default:
				h.log.Error(err.Error())
				err = httputils.RespondWith500(w)
			}

			if err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		resp := &getReviewResponse{
			Reviews: make([]respReview, 0, len(reviews)),
		}

		for _, review := range reviews {
			resp.Reviews = append(resp.Reviews, respReview{
				ReviewID:       review.ID,
				TeacherID:      review.TeacherID,
				SkillID:        review.SkillID,
				CategoryID:     review.CategoryID,
				Rate:           review.Rate,
				Comment:        review.Comment,
				StudentID:      review.StudentID,
				StudentEmail:   review.StudentData.Email,
				StudentName:    review.StudentData.Name,
				StudentSurname: review.StudentData.Surname,
				StudentAvatar:  review.StudentData.Avatar,
			})
		}

		respondErr := httputils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			h.log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

type getReviewResponse struct {
	Reviews []respReview `json:"reviews"`
}

type respReview struct {
	ReviewID       int    `json:"review_id"   example:"1"`
	TeacherID      int    `json:"teacher_id"  example:"1"`
	SkillID        int    `json:"skill_id"    example:"1"`
	CategoryID     int    `json:"category_id" example:"1"`
	Rate           int    `json:"rate"        example:"5"`
	Comment        string `json:"comment"     example:"This is a comment"`
	StudentID      int    `json:"student_id"      example:"1"`
	StudentEmail   string `json:"student_email"   example:"qwerty@example.com"`
	StudentName    string `json:"student_name"    example:"John"`
	StudentSurname string `json:"student_surname" example:"Smith"`
	StudentAvatar  string `json:"student_avatar"  example:"uuid.png"`
}
