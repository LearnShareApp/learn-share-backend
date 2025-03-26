package add_review

import (
	"encoding/json"
	"errors"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
	"net/http"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"go.uber.org/zap"
)

const (
	Route = ""
)

// MakeHandler returns http.HandlerFunc
// @Summary Create review
// @Description Create review if authorized user (student) had lesson with this teacher and this category
// @Tags reviews
// @Accept json
// @Produce json
// @Param request body request true "Review data"
// @Success 201
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 403 {object} httputils.ErrorStruct
// @Failure 404 {object} httputils.ErrorStruct
// @Failure 409 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /review [post]
// @Security     BearerAuth
func MakeHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userId := r.Context().Value(jwt.UserIDKey).(int)
		if userId == 0 {
			log.Error("id was missed in context")
			if err := httputils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		var req request

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			if err = httputils.RespondWith400(w, "failed to decode body"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		if req.TeacherId == 0 || req.CategoryId == 0 || req.Rate == 0 || req.Comment == "" {
			if err := httputils.RespondWith400(w, "teacher_id, category_id, rate or comment are empty"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		if req.Rate > 5 {
			if err := httputils.RespondWith400(w, "rate must be less or equal than 5"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		review := entities.Review{
			TeacherId:  req.TeacherId,
			StudentId:  userId,
			CategoryId: req.CategoryId,
			Rate:       req.Rate,
			Comment:    req.Comment,
		}

		err := s.Do(r.Context(), &review)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorUserNotFound):
				err = httputils.RespondWith401(w, serviceErrors.ErrorUserNotFound.Error())
			case errors.Is(err, serviceErrors.ErrorTeacherNotFound):
				err = httputils.RespondWith404(w, serviceErrors.ErrorTeacherNotFound.Error())
			case errors.Is(err, serviceErrors.ErrorStudentAndTeacherSame):
				err = httputils.RespondWith400(w, "teacher can not review himself")
			case errors.Is(err, serviceErrors.ErrorCategoryNotFound):
				err = httputils.RespondWith404(w, serviceErrors.ErrorCategoryNotFound.Error())
			case errors.Is(err, serviceErrors.ErrorSkillUnregistered):
				err = httputils.RespondWith404(w, serviceErrors.ErrorSkillUnregistered.Error())
			case errors.Is(err, serviceErrors.ErrorFinishedLessonNotFound):
				err = httputils.RespondWith403(w, "You have not finished lesson with this teacher and this category")
			case errors.Is(err, serviceErrors.ErrorReviewExists):
				err = httputils.RespondWith409(w, serviceErrors.ErrorReviewExists.Error())
			default:
				log.Error(err.Error())
				err = httputils.RespondWith500(w)
			}

			if err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}
		respondErr := httputils.SuccessRespondWith201(w, struct{}{})
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}
