package get_reviews

import (
	"errors"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"go.uber.org/zap"
	"net/http"
)

const (
	Route = "/{id}/reviews"
)

// MakeHandler returns http.HandlerFunc
// @Summary Get reviews by teacher's ID
// @Description Get all reviews by teacher's ID
// @Tags teachers
// @Accept json
// @Produce json
// @Param id path int true "Teacher's ID"
// @Success 200 {object} response
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 404 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /teachers/{id}/reviews [get]
func MakeHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teacherId, err := httputils.GetIntParamFromRequestPath(r, "id")
		if err != nil {
			if err := httputils.RespondWith400(w, "missed or not-number {id} param in url path"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		reviews, err := s.Do(r.Context(), teacherId)
		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorTeacherNotFound):
				err = httputils.RespondWith404(w, serviceErrors.ErrorTeacherNotFound.Error())
			default:
				log.Error(err.Error())
				err = httputils.RespondWith500(w)
			}

			if err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		resp := response{
			Reviews: make([]review, 0, len(reviews)),
		}

		for _, rev := range reviews {
			resp.Reviews = append(resp.Reviews, review{
				ReviewId:   rev.Id,
				TeacherId:  rev.TeacherId,
				SkillId:    rev.SkillId,
				CategoryId: rev.CategoryId,
				Rate:       rev.Rate,
				Comment:    rev.Comment,

				StudentId:      rev.StudentId,
				StudentEmail:   rev.StudentData.Email,
				StudentName:    rev.StudentData.Name,
				StudentSurname: rev.StudentData.Surname,
				StudentAvatar:  rev.StudentData.Avatar,
			})

		}

		respondErr := httputils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}
