package registration

import (
	"encoding/json"
	"github.com/LearnShareApp/learn-share-backend/internal/jsonutils"
	"go.uber.org/zap"
	"net/http"
)

func MakeHandler(s *Service, log *zap.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode body", zap.Error(err))
			if err = jsonutils.RespondWith400(w, "failed to decode body"); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}
	}
}
