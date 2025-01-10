package get

import (
	"github.com/LearnShareApp/learn-share-backend/internal/jsonutils"
	"go.uber.org/zap"
	"net/http"
)

func MakeHandler(s *Service, log *zap.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		categories, err := s.Do(r.Context())
		if err != nil {
			log.Error(err.Error())
			if err = jsonutils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		resp := &response{
			Categories: make([]category, 0, len(categories)),
		}
		for _, c := range categories {
			resp.Categories = append(resp.Categories, category{
				Name:   c.Name,
				MinAge: int(c.MinAge),
			})
		}

		respondErr := jsonutils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}
