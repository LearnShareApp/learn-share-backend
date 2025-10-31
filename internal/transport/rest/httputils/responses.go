package httputils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type ErrorStruct struct {
	Error string `json:"error"`
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		respondErr := RespondWithError(w,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
		if respondErr != nil {
			return err
		}

		return fmt.Errorf("failed to marshall payload: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if _, err = w.Write(response); err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}

func RespondWithImage(w http.ResponseWriter, code int, reader io.Reader, imgExtension string) error {
	if imgExtension != "png" && imgExtension != "jpg" {
		respondErr := RespondWithError(w,
			http.StatusInternalServerError,
			"Unsupported image extension")
		if respondErr != nil {
			return respondErr
		}

		return fmt.Errorf("unsupported image extension: %s", imgExtension)
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "image/"+imgExtension)
	w.Header().Set("Cache-Control", "public, max-age=3600") // cache for 1 hour

	if _, err := io.Copy(w, reader); err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}

func RespondWithError(w http.ResponseWriter, code int, message string) error {
	return RespondWithJSON(w,
		code,
		ErrorStruct{
			Error: message,
		})
}

func RespondWith400(w http.ResponseWriter, message string, log *zap.Logger) {
	if err := RespondWithError(w,
		http.StatusBadRequest,
		message); err != nil {
		log.Error("response error", zap.Error(err))
	}
}

func RespondWith401(w http.ResponseWriter, message string, log *zap.Logger) {
	if err := RespondWithError(w,
		http.StatusUnauthorized,
		message); err != nil {
		log.Error("response error", zap.Error(err))
	}
}

func RespondWith403(w http.ResponseWriter, message string, log *zap.Logger) {
	if err := RespondWithError(w,
		http.StatusForbidden,
		message); err != nil {
		log.Error("response error", zap.Error(err))
	}
}

func RespondWith404(w http.ResponseWriter, message string, log *zap.Logger) {
	if err := RespondWithError(w,
		http.StatusNotFound,
		message); err != nil {
		log.Error("response error", zap.Error(err))
	}
}

func RespondWith409(w http.ResponseWriter, message string, log *zap.Logger) {
	if err := RespondWithError(w,
		http.StatusConflict,
		message); err != nil {
		log.Error("response error", zap.Error(err))
	}
}

func RespondWith500(w http.ResponseWriter, log *zap.Logger) {
	if err := RespondWithError(w,
		http.StatusInternalServerError,
		http.StatusText(http.StatusInternalServerError)); err != nil {
		log.Error("response error", zap.Error(err))
	}
}

func SuccessRespondWith200(w http.ResponseWriter, payload interface{}, log *zap.Logger) {
	if err := RespondWithJSON(w,
		http.StatusOK,
		payload); err != nil {
		log.Error("response error", zap.Error(err))
	}
}

func SuccessRespondWith201(w http.ResponseWriter, payload interface{}, log *zap.Logger) {
	if err := RespondWithJSON(w,
		http.StatusCreated,
		payload); err != nil {
		log.Error("response error", zap.Error(err))
	}
}
