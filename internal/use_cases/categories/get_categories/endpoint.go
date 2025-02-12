package get_categories

import (
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"go.uber.org/zap"
	"net/http"
)

const Route = "/categories"

// MakeHandler returns http.HandlerFunc
// @Summary Get categories
// @Description Get list of all categories
// @Tags categories
// @Produce json
// @Success 200 {object} response
// @Failure 500 {object} httputils.ErrorStruct
// @Router /categories [get]
func MakeHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		categories, err := s.Do(r.Context())
		if err != nil {
			log.Error(err.Error())
			if err = httputils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		resp := &response{
			Categories: make([]category, 0, len(categories)),
		}
		for _, c := range categories {
			resp.Categories = append(resp.Categories, category{
				Id:     int(c.Id),
				Name:   c.Name,
				MinAge: int(c.MinAge),
			})
		}

		respondErr := httputils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}
