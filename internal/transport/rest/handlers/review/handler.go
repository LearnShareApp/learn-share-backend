package review

import (
	"context"
	"net/http"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ReviewService interface {
	CreateReview(ctx context.Context, review *entities.Review) error
	GetReviews(ctx context.Context, teacherID int) ([]*entities.Review, error)
}

type ReviewHandlers struct {
	reviewService ReviewService
	log           *zap.Logger
}

func NewReviewHandlers(reviewService ReviewService, log *zap.Logger) *ReviewHandlers {
	return &ReviewHandlers{
		reviewService: reviewService,
		log:           log,
	}
}

func (h *ReviewHandlers) SetupReviewRoutes(router *chi.Mux, authMiddleware func(http.Handler) http.Handler) {
	router.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post(createRoute, h.CreateReview())
	})

	router.Get(getListRoute, h.GetReviewList())
}
