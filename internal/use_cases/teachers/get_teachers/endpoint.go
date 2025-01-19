package get_teachers

import (
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	"github.com/LearnShareApp/learn-share-backend/internal/jsonutils"
	"github.com/LearnShareApp/learn-share-backend/internal/service/jwt"
	"go.uber.org/zap"
	"net/http"
)

const (
	Route = ""
)

// MakeHandler returns http.HandlerFunc which handle get teachers
// @Summary Get full teachers data
// @Description Get full teachers data (their user data, teacher data and skills)
// @Tags teachers
// @Produce json
// @Success 200 {object} response
// @Failure 401 {object} jsonutils.ErrorStruct
// @Failure 500 {object} jsonutils.ErrorStruct
// @Router /teachers [get]
// @Security     BearerAuth
func MakeHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.Context().Value(jwt.UserIDKey).(int)
		if id == 0 {
			log.Error("id was missed in context")
			if err := jsonutils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		// TODO filters: by category_id, by personal teachers (student - teacher)
		teachers, err := s.Do(r.Context())

		if err != nil {
			log.Error(err.Error())
			if err = jsonutils.RespondWith500(w); err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		resp := mappingResponse(teachers)

		respondErr := jsonutils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

func mappingResponse(users []entities.User) *response {
	resp := &response{
		Teachers: make([]teacher, 0, len(users)),
	}

	for i := range users {
		skills := make([]skill, len(users[i].TeacherData.Skills))
		// Transform skills slice in one go
		for j, sk := range users[i].TeacherData.Skills {
			skills[j] = skill{
				SkillId:       sk.Id,
				CategoryId:    sk.CategoryId,
				CategoryName:  sk.CategoryName,
				VideoCardLink: sk.VideoCardLink,
				About:         sk.About,
				Rate:          sk.Rate,
			}
		}

		resp.Teachers = append(resp.Teachers, teacher{
			TeacherId:        users[i].TeacherData.Id,
			UserId:           users[i].Id,
			Email:            users[i].Email,
			Name:             users[i].Name,
			Surname:          users[i].Surname,
			RegistrationDate: users[i].RegistrationDate,
			Birthdate:        users[i].Birthdate,
			Skills:           skills,
		})
	}

	return resp
}
