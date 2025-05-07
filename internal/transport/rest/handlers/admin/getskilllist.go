package admin

import (
	"net/http"
	"strconv"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
	"go.uber.org/zap"
)

const getSkillListRoute = "/skills"

// GetSkillList returns http.HandlerFunc
// @Summary get skill's list (and their teachers sort data)
// @Description returns the list of skills and have one flag unactive, if it's true, returns only anactive else => all (+ skill's teachers sort data as addtional list)
// @Tags admin
// @Produce json
// @Param unactive query boolean false "flag for only unactive"
// @Success 200 {object} getSkillListResponse
// @Failure 401 {object} httputils.ErrorStruct
// @Failure 403 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /admin/skills [get]
// @Security     BearerAuth
func (h *AdminHandlers) GetSkillList() http.HandlerFunc {
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

		// get query flag
		isUnactive := r.URL.Query().Get("unactive")
		isUnactiveBool, err := strconv.ParseBool(isUnactive)
		if err != nil {
			isUnactiveBool = false
		}

		isAdmin, err := h.service.CheckUserOnAdminByID(r.Context(), userID)

		if err != nil {
			h.log.Error("failed to check user on admin", zap.Error(err))
			if err := httputils.RespondWith500(w); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}

			return
		}

		if !isAdmin {
			if err := httputils.RespondWith403(w, serviceErrors.ErrorNotAdmin.Error()); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}

			return
		}

		var skills []entities.Skill

		if isUnactiveBool {
			skills, err = h.service.GetUnactiveSkillList(r.Context())
		} else {
			skills, err = h.service.GetSkillList(r.Context())
		}

		if err != nil {
			h.log.Error(err.Error())
			err = httputils.RespondWith500(w)

			if err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}

			return
		}

		teacherIDs := make([]int, 0)
		for i := range skills {
			teacherIDs = append(teacherIDs, skills[i].TeacherID)
		}

		teachersUserData, err := h.service.GetTeacherShortDataListByIDs(r.Context(), teacherIDs)
		if err != nil {
			h.log.Error(err.Error())
			err = httputils.RespondWith500(w)

			if err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}

			return
		}

		resp := getSkillListResponse{
			Skills:   make([]respSkill, 0, len(skills)),
			Teachers: make([]respTeacherShortData, 0, len(teachersUserData)),
		}

		for i := range skills {
			resp.Skills = append(resp.Skills, respSkill{
				SkillID:       skills[i].ID,
				TeacherID:     skills[i].TeacherID,
				CategoryID:    skills[i].CategoryID,
				VideoCardLink: skills[i].VideoCardLink,
				About:         skills[i].About,
				Rate:          skills[i].Rate,
				ReviewsCount:  skills[i].ReviewsCount,
				IsActive:      skills[i].IsActive,
			})
		}

		for i := range teachersUserData {
			resp.Teachers = append(resp.Teachers, respTeacherShortData{
				TeacherID: teachersUserData[i].TeacherData.ID,
				Name:      teachersUserData[i].Name,
				Surname:   teachersUserData[i].Surname,
				Avatar:    teachersUserData[i].Avatar,
			})
		}

		if err = httputils.SuccessRespondWith200(w, resp); err != nil {
			h.log.Error("failed to send response", zap.Error(err))
		}
	}
}

type getSkillListResponse struct {
	Skills   []respSkill            `json:"skills"`
	Teachers []respTeacherShortData `json:"teachers"`
}

type respSkill struct {
	SkillID       int     `json:"skill_id"        example:"1"`
	TeacherID     int     `json:"teacher_id"      example:"1"`
	CategoryID    int     `json:"category_id"     example:"1"`
	VideoCardLink string  `json:"video_card_link" example:"https://youtu.be/HIcSWuKMwOw?si=FtxN1QJU9ZWnXy85"`
	About         string  `json:"about"           example:"about me..."`
	Rate          float32 `json:"rate"            example:"5"`
	ReviewsCount  int     `json:"reviews_count"   example:"1"`
	IsActive      bool    `json:"is_active"       example:"false"`
}

type respTeacherShortData struct {
	TeacherID int    `json:"teacher_id"      example:"1"`
	Name      string `json:"name"            example:"John"`
	Surname   string `json:"surname"         example:"Smith"`
	Avatar    string `json:"avatar"          example:"uuid.png"`
}
