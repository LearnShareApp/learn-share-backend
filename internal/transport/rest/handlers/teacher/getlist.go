package teacher

import (
	"net/http"
	"strconv"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
	"go.uber.org/zap"
)

const (
	getTeacherListRoute = "/"
)

// GetTeacherList returns http.HandlerFunc which handle get teachers
// @Summary Get full teachers data
// @Description Get full teachers data (their user data, teacher data and skills)
// @Tags teachers
// @Produce json
// @Param is_mine query boolean false "Filter my teachers"
// @Param category query string false "Filter category"
// @Success 200 {object} getTeachersResponse
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /teachers [get]
// @Security     BearerAuth
func (h *TeacherHandlers) GetTeacherList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDValue := r.Context().Value(jwt.UserIDKey)
		id, ok := userIDValue.(int)
		if !ok || id == 0 {
			h.log.Error("invalid or missing user ID in context", zap.Any("value", userIDValue))
			if err := httputils.RespondWith500(w); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		// get query as filters
		isMyTeachers := r.URL.Query().Get("is_mine")
		isMyTeachersBool, err := strconv.ParseBool(isMyTeachers)
		if err != nil {
			isMyTeachersBool = false
		}

		category := r.URL.Query().Get("category")
		isFilterByCategory := true
		if category == "" {
			isFilterByCategory = false
		}

		teachers, err := h.teacherService.GetTeacherList(r.Context(), id, isMyTeachersBool, category, isFilterByCategory)

		if err != nil {
			h.log.Error(err.Error())
			if err = httputils.RespondWith500(w); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		resp := mappingResponse(teachers)

		respondErr := httputils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			h.log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

func mappingResponse(users []entities.User) *getTeachersResponse {
	resp := &getTeachersResponse{
		Teachers: make([]getTeacherResponse, 0, len(users)),
	}

	for i := range users {
		skills := make([]respSkill, len(users[i].TeacherData.Skills))
		// Transform skills slice in one go
		for j, sk := range users[i].TeacherData.Skills {
			skills[j] = respSkill{
				SkillID:       sk.Id,
				CategoryID:    sk.CategoryId,
				CategoryName:  sk.CategoryName,
				VideoCardLink: sk.VideoCardLink,
				About:         sk.About,
				Rate:          sk.Rate,
				ReviewsCount:  sk.ReviewsCount,
			}
		}

		resp.Teachers = append(resp.Teachers, getTeacherResponse{
			TeacherID:          users[i].TeacherData.Id,
			UserID:             users[i].Id,
			Email:              users[i].Email,
			Name:               users[i].Name,
			Surname:            users[i].Surname,
			RegistrationDate:   users[i].RegistrationDate,
			Birthdate:          users[i].Birthdate,
			Avatar:             users[i].Avatar,
			FinishedLessons:    users[i].TeacherData.TeacherStat.CountOfFinishedLesson,
			CountOfStudents:    users[i].TeacherData.TeacherStat.CountOfStudents,
			CommonRate:         users[i].TeacherData.Rate,
			CommonReviewsCount: users[i].TeacherData.ReviewsCount,
			Skills:             skills,
		})
	}

	return resp
}

type getTeachersResponse struct {
	Teachers []getTeacherResponse `json:"teachers"`
}
