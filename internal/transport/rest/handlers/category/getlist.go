package category

import (
	"net/http"

	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"go.uber.org/zap"
)

const Route = "/categories"

// GetCategoryList returns http.HandlerFunc
// @Summary Get categories
// @Description Get list of all categories
// @Tags categories
// @Produce json
// @Success 200 {object} getCategoriesResponse
// @Failure 500 {object} httputils.ErrorStruct
// @Router /categories [get]
func (h *CategoryHandlers) GetCategoryList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, err := h.categoryService.GetCategories(r.Context())
		if err != nil {
			h.log.Error(err.Error())
			if err = httputils.RespondWith500(w); err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		resp := &getCategoriesResponse{
			Categories: make([]respCategory, 0, len(categories)),
		}
		for _, c := range categories {
			resp.Categories = append(resp.Categories, respCategory{
				ID:     c.Id,
				Name:   c.Name,
				MinAge: c.MinAge,
			})
		}

		respondErr := httputils.SuccessRespondWith200(w, resp)
		if respondErr != nil {
			h.log.Error("failed to send response", zap.Error(respondErr))
		}
	}
}

// @Description get categories getCategoriesResponse.
type getCategoriesResponse struct {
	Categories []respCategory `json:"categories"`
}

// @Description data of respCategory.
type respCategory struct {
	ID     int    `json:"id"      example:"1"`
	Name   string `json:"name"    example:"Programing"`
	MinAge int    `json:"min_age" example:"12"`
}
