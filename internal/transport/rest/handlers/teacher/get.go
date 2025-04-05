package teacher

import (
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
)

const (
	getTeacherProtectedRoute = "/"
	getTeacherPublicRoute    = "/{id}"
)

// GetTeacherProtected returns http.HandlerFunc which handle get teacher, get user id from token
// @Summary Get teacher data
// @Description Get all info about teacher (user info + teacher + his skills) by user id in token
// @Tags teachers
// @Produce json
// @Success 200 {object} getTeacherResponse
// @Failure 403 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /teacher [get]
// @Security     BearerAuth
func (h *TeacherHandlers) GetTeacherProtected() http.HandlerFunc {
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

		teacher, err := h.teacherService.GetTeacherByUserID(r.Context(), userID)

		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorTeacherNotFound):
				err = httputils.RespondWith403(w, serviceErrors.ErrorUserIsNotTeacher.Error())
			default:
				h.log.Error(err.Error())
				err = httputils.RespondWith500(w)
			}

			if err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}

			return
		}

		teacherAllData, err := h.teacherService.GetTeacher(r.Context(), teacher)

		if err != nil {
			coveringErrors(w, h.log, err)
			return
		}

		resp := mappingToResponse(teacherAllData)

		respondErr := httputils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			h.log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

// GetTeacherPublic returns http.HandlerFunc which handle get teacher, get user id from http param
// @Summary Get teacher data
// @Description Get all info about teacher (user info + teacher + his skills) by his TeacherID in route (/api/teachers/{id})
// @Tags teachers
// @Produce json
// @Param id path int true "Teacher's ID"
// @Success 200 {object} getTeacherResponse
// @Failure 404 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /teachers/{id} [get]
func (h *TeacherHandlers) GetTeacherPublic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teacherID, err := httputils.GetIntParamFromRequestPath(r, "id")
		if err != nil {
			if err := httputils.RespondWith400(w, "missed or not-number {id} param in url path"); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		teacher, err := h.teacherService.GetTeacherByID(r.Context(), teacherID)

		if err != nil {
			coveringErrors(w, h.log, err)
			return
		}

		teacherAllData, err := h.teacherService.GetTeacher(r.Context(), teacher)

		if err != nil {
			coveringErrors(w, h.log, err)
			return
		}

		resp := mappingToResponse(teacherAllData)

		respondErr := httputils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			h.log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

func coveringErrors(w http.ResponseWriter, log *zap.Logger, err error) {
	switch {
	case errors.Is(err, serviceErrors.ErrorTeacherNotFound):
		err = httputils.RespondWith404(w, err.Error())
	default:
		log.Error(err.Error())
		err = httputils.RespondWith500(w)
	}

	if err != nil {
		log.Error("failed to send response", zap.Error(err))
	}
}

func mappingToResponse(user *entities.User) *getTeacherResponse {
	resp := getTeacherResponse{
		TeacherID:          user.TeacherData.Id,
		UserID:             user.Id,
		Email:              user.Email,
		Name:               user.Name,
		Surname:            user.Surname,
		RegistrationDate:   user.RegistrationDate,
		Birthdate:          user.Birthdate,
		Avatar:             user.Avatar,
		FinishedLessons:    user.TeacherData.TeacherStat.CountOfFinishedLesson,
		CountOfStudents:    user.TeacherData.TeacherStat.CountOfStudents,
		CommonRate:         user.TeacherData.Rate,
		CommonReviewsCount: user.TeacherData.ReviewsCount,

		Skills: make([]respSkill, 0, len(user.TeacherData.Skills)),
	}

	// remap entity respSkill to getTeacherResponse respSkill-type
	for _, sk := range user.TeacherData.Skills {
		resp.Skills = append(resp.Skills, respSkill{
			SkillID:       sk.Id,
			CategoryID:    sk.CategoryId,
			CategoryName:  sk.CategoryName,
			VideoCardLink: sk.VideoCardLink,
			About:         sk.About,
			Rate:          sk.Rate,
			ReviewsCount:  sk.ReviewsCount,
		})
	}

	return &resp
}

type getTeacherResponse struct {
	TeacherID          int         `json:"teacher_id"           example:"1"`
	UserID             int         `json:"user_id"              example:"1"`
	Email              string      `json:"email"                example:"qwerty@example.com"`
	Name               string      `json:"name"                 example:"John"`
	Surname            string      `json:"surname"              example:"Smith"`
	RegistrationDate   time.Time   `json:"registration_date"    example:"2022-09-09T10:10:10+09:00"`
	Birthdate          time.Time   `json:"birthdate"            example:"2002-09-09T10:10:10+09:00"`
	Avatar             string      `json:"avatar"               example:"uuid.png"`
	FinishedLessons    int         `json:"finished_lessons"     example:"0"`
	CountOfStudents    int         `json:"count_of_students"    example:"0"`
	CommonRate         float32     `json:"common_rate"          example:"0"`
	CommonReviewsCount int         `json:"common_reviews_count" example:"0"`
	Skills             []respSkill `json:"skills"`
}

type respSkill struct {
	SkillID       int     `json:"skill_id"        example:"1"`
	CategoryID    int     `json:"category_id"     example:"1"`
	CategoryName  string  `json:"category_name"   example:"Category"`
	VideoCardLink string  `json:"video_card_link" example:"https://youtu.be/HIcSWuKMwOw?si=FtxN1QJU9ZWnXy85"`
	About         string  `json:"about"           example:"about me..."`
	Rate          float32 `json:"rate"            example:"5"`
	ReviewsCount  int     `json:"reviews_count"   example:"1"`
}
