package image

import (
	"context"
	"io"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ImageService interface {
	GetImage(ctx context.Context, filename string) (io.Reader, error)
}

type ImageHandlers struct {
	imageService ImageService
	log          *zap.Logger
}

func NewImageHandlers(imageService ImageService, log *zap.Logger) *ImageHandlers {
	return &ImageHandlers{
		imageService: imageService,
		log:          log,
	}
}

func (h *ImageHandlers) SetupImageRoutes(router *chi.Mux) {
	router.Get(GetRoute, h.GetImage())
}
