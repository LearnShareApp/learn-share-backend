package get_image

import (
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"slices"
	"strings"

	serviceErrors "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/internal/jsonutils"
	"go.uber.org/zap"
)

const Route = "/image"

// MakeHandler returns http.HandlerFunc
// @Summary Get image
// @Description Get image by filename
// @Tags image
// @Produce image/*
// @Param filename query string true "filename"
// @Success 200 {file} binary "Image file"
// @Failure 400 {object} jsonutils.ErrorStruct
// @Failure 404 {object} jsonutils.ErrorStruct
// @Failure 500 {object} jsonutils.ErrorStruct
// @Router /image [get]
func MakeHandler(s *Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Query().Get("filename")
		if filename == "" {
			if err := jsonutils.RespondWith400(w, "missing filename parameter"); err != nil {
				log.Error("failed to respond with 400", zap.Error(err))
			}
			return
		}

		extension := strings.TrimPrefix(filepath.Ext(filename), ".")
		if extension == "" {
			if err := jsonutils.RespondWith400(w, "missing extension in filename"); err != nil {
				log.Error("failed to respond with 400", zap.Error(err))
			}
			return
		}

		if !slices.Contains(SupportedExtension, extension) {
			if err := jsonutils.RespondWith400(w, "unsupported extension"); err != nil {
				log.Error("failed to respond with 400", zap.Error(err))
			}
			return
		}

		reader, err := s.Do(r.Context(), filename)
		defer func() {
			if closer, ok := reader.(io.Closer); ok {
				if err := closer.Close(); err != nil {
					log.Error("failed to close reader", zap.Error(err))
				}
			}
		}()

		if err != nil {
			switch {
			case errors.Is(err, serviceErrors.ErrorImageNotFound):
				err = jsonutils.RespondWith404(w, err.Error())
			case errors.Is(err, serviceErrors.ErrorIncorrectFileFormat):
				err = jsonutils.RespondWith400(w, err.Error())
			default:
				log.Error(err.Error())
				err = jsonutils.RespondWith500(w)
			}

			if err != nil {
				log.Error("failed to send response", zap.Error(err))
			}
			return
		}

		if err := jsonutils.RespondWithImage(w, http.StatusOK, reader, extension); err != nil {
			log.Error("failed to send response", zap.Error(err))
		}
	}
}
