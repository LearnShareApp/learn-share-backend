package category

import (
	"context"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type CategoryService interface {
	GetCategories(ctx context.Context) ([]*entities.Category, error)
}

type CategoryHandlers struct {
	categoryService CategoryService
	log             *zap.Logger
}

func NewCategoryHandlers(categoryService CategoryService, log *zap.Logger) *CategoryHandlers {
	return &CategoryHandlers{
		categoryService: categoryService,
		log:             log,
	}
}

func (h *CategoryHandlers) SetupCategoryRoutes(router *chi.Mux) {
	router.Get(Route, h.GetCategoryList())
}
