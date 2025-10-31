package review

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/httputils"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
	"go.uber.org/zap"
)

const (
	createRoute = "/review"
)

// CreateReview returns http.HandlerFunc
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
func (h *ReviewHandlers) CreateReview() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDValue := r.Context().Value(jwt.UserIDKey)
		userID, ok := userIDValue.(int)
		if !ok || userID == 0 {
			h.log.Error("invalid or missing user ID in context", zap.Any("value", userIDValue))
			httputils.RespondWith500(w, h.log)

			return
		}

		var req addReviewRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httputils.RespondWith400(w, "failed to decode body", h.log)

			return
		}

		if req.TeacherID == 0 || req.CategoryID == 0 || req.Rate == 0 || req.Comment == "" {
			httputils.RespondWith400(w, "teacher_id, category_id, rate or comment are empty", h.log)

			return
		}

		if req.Rate > 5 {
			httputils.RespondWith400(w, "rate must be less or equal than 5", h.log)

			return
		}

		review := &entities.Review{
			TeacherID:  req.TeacherID,
			StudentID:  userID,
			CategoryID: req.CategoryID,
			Rate:       req.Rate,
			Comment:    req.Comment,
		}

		err := h.reviewService.CreateReview(r.Context(), review)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				httputils.RespondWith401(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorTeacherNotFound):
				httputils.RespondWith404(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorStudentAndTeacherSame):
				httputils.RespondWith403(w, "teacher can not review himself", h.log)
			case errors.Is(err, serviceErrors.ErrorCategoryNotFound):
				httputils.RespondWith404(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorSkillUnregistered):
				httputils.RespondWith404(w, err.Error(), h.log)
			case errors.Is(err, serviceErrors.ErrorFinishedLessonNotFound):
				httputils.RespondWith403(w, "You have not finished lesson with this teacher and this category", h.log)
			case errors.Is(err, serviceErrors.ErrorReviewExists):
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

type addReviewRequest struct {
	TeacherID  int    `json:"teacher_id"  example:"1"            binding:"required"`
	CategoryID int    `json:"category_id" example:"1"            binding:"required"`
	Rate       int    `json:"rate"        example:"1"            binding:"required"`
	Comment    string `json:"comment"     example:"some comment" binding:"required"`
}
