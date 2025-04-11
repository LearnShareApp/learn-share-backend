package review

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
	"go.uber.org/zap"
)

const (
	addRoute = "/review"
)

// AddReview returns http.HandlerFunc
// @Summary Create review
// @Description Create review if authorized user (student) had lesson with this teacher and this category
// @Tags reviews
// @Accept json
// @Produce json
// @Param addReviewRequest body addReviewRequest true "Review data"
// @Success 201
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 403 {object} httputils.ErrorStruct
// @Failure 404 {object} httputils.ErrorStruct
// @Failure 409 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /review [post]
// @Security     BearerAuth
func (h *ReviewHandlers) AddReview() http.HandlerFunc {
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

		var req addReviewRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			if err = httputils.RespondWith400(w, "failed to decode body"); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		if req.TeacherID == 0 || req.CategoryID == 0 || req.Rate == 0 || req.Comment == "" {
			if err := httputils.RespondWith400(w, "teacher_id, category_id, rate or comment are empty"); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		if req.Rate > 5 {
			if err := httputils.RespondWith400(w, "rate must be less or equal than 5"); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		review := &entities.Review{
			TeacherID:  req.TeacherID,
			StudentID:  userID,
			CategoryID: req.CategoryID,
			Rate:       req.Rate,
			Comment:    req.Comment,
		}

		err := h.reviewService.AddReview(r.Context(), review)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				err = httputils.RespondWith401(w, err.Error())
			case errors.Is(err, serviceErrors.ErrorTeacherNotFound):
				err = httputils.RespondWith404(w, err.Error())
			case errors.Is(err, serviceErrors.ErrorStudentAndTeacherSame):
				err = httputils.RespondWith403(w, "teacher can not review himself")
			case errors.Is(err, serviceErrors.ErrorCategoryNotFound):
				err = httputils.RespondWith404(w, err.Error())
			case errors.Is(err, serviceErrors.ErrorSkillUnregistered):
				err = httputils.RespondWith404(w, err.Error())
			case errors.Is(err, serviceErrors.ErrorFinishedLessonNotFound):
				err = httputils.RespondWith403(w, "You have not finished lesson with this teacher and this category")
			case errors.Is(err, serviceErrors.ErrorReviewExists):
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

type addReviewRequest struct {
	TeacherID  int    `json:"teacher_id"  example:"1"            binding:"required"`
	CategoryID int    `json:"category_id" example:"1"            binding:"required"`
	Rate       int    `json:"rate"        example:"1"            binding:"required"`
	Comment    string `json:"comment"     example:"some comment" binding:"required"`
}
