package image

import (
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"slices"
	"strings"

	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"go.uber.org/zap"
)

const GetRoute = "/image"

var SupportedExtension = []string{"png", "jpg", "jpeg"}

// GetImage returns http.HandlerFunc
// @Summary Get image
// @Description Get image by filename
// @Tags image
// @Produce image/*
// @Param filename query string true "filename"
// @Success 200 {file} binary "Image file"
// @Failure 400 {object} httputils.ErrorStruct
// @Failure 404 {object} httputils.ErrorStruct
// @Failure 500 {object} httputils.ErrorStruct
// @Router /image [get]
func (h *ImageHandlers) GetImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Query().Get("filename")
		if filename == "" {
			if err := httputils.RespondWith400(w, "missing filename parameter"); err != nil {
				h.log.Error("failed to respond with 400", zap.Error(err))
			}

			return
		}

		extension := strings.TrimPrefix(filepath.Ext(filename), ".")
		if extension == "" {
			if err := httputils.RespondWith400(w, "missing extension in filename"); err != nil {
				h.log.Error("failed to respond with 400", zap.Error(err))
			}

			return
		}

		if !slices.Contains(SupportedExtension, extension) {
			if err := httputils.RespondWith400(w, "unsupported extension"); err != nil {
				h.log.Error("failed to respond with 400", zap.Error(err))
			}
			return
		}

		reader, err := h.imageService.GetImage(r.Context(), filename)
		defer func() {
			if closer, ok := reader.(io.Closer); ok {
				if err := closer.Close(); err != nil {
					h.log.Error("failed to close reader", zap.Error(err))
				}
			}
		}()

		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorImageNotFound):
				err = httputils.RespondWith404(w, err.Error())
			case errors.Is(err, serviceErrors.ErrorIncorrectFileFormat):
				err = httputils.RespondWith400(w, err.Error())
			default:
				h.log.Error(err.Error())
				err = httputils.RespondWith500(w)
			}

			if err != nil {
				h.log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		if err := httputils.RespondWithImage(w, http.StatusOK, reader, extension); err != nil {
			h.log.Error("failed to send response", zap.Error(err))
		}
	}
}
