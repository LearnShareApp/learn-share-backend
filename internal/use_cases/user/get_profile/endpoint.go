package get_profile

import (
	"github.com/LearnShareApp/learn-share-backend/internal/jsonutils"
	"github.com/LearnShareApp/learn-share-backend/internal/service/jwt"
	"go.uber.org/zap"
	"net/http"
)

const Route = "/profile"

// MakeHandler returns http.HandlerFunc
// @Summary Get user profile
// @Description Get info about user by jwt token (in Authorization enter: Bearer <your_jwt_token>)
// @Tags user
// @Produce json
// @Success 200 {object} response
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /user/profile [get]
// @Security     BearerAuth
func MakeHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.Context().Value(jwt.UserIDKey).(int64)
		if id == 0 {
			log.Error("id was missed in context")
			if err := jsonutils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		user, err := s.Do(r.Context(), id)
		if err != nil {
			// Вроде как нет смысла обрабатывать кейс когда пользователь не найден по id т. к. по хорошему в токене,
			// который выпустили мы не может быть несуществующий пользователь
			// что и есть 500-я

			log.Error(err.Error())
			if err = jsonutils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		resp := response{
			Email:            user.Email,
			Name:             user.Name,
			Surname:          user.Surname,
			RegistrationDate: user.RegistrationDate,
			Birthdate:        user.Birthdate,
		}

		respondErr := jsonutils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}
