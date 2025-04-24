package complaint

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	"github.com/go-chi/chi/v5"
)

const (
	complaintRoute = "/complaint"
)

type ComplaintService interface {
	CreateComplaint(ctx context.Context, complaint *entities.Complaint) error
}

type ComplaintHandlers struct {
	service ComplaintService
	log     *zap.Logger
}

func NewComplaintHandlers(complaintService ComplaintService, log *zap.Logger) *ComplaintHandlers {
	return &ComplaintHandlers{
		service: complaintService,
		log:     log,
	}
}

func (h *ComplaintHandlers) SetupComplaintRoutes(router *chi.Mux, authMiddleware func(http.Handler) http.Handler) {
	complaintRouter := chi.NewRouter()

	complaintRouter.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		r.Post(createRoute, h.CreateComplaint())
	})

	router.Mount(complaintRoute, complaintRouter)
}
