package httputils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func RespondWith400(w http.ResponseWriter, message string) error {
	return RespondWithError(w,
		http.StatusBadRequest,
		message)
}

func RespondWith401(w http.ResponseWriter, message string) error {
	return RespondWithError(w,
		http.StatusUnauthorized,
		message)
}

func RespondWith403(w http.ResponseWriter, message string) error {
	return RespondWithError(w,
		http.StatusForbidden,
		message)
}

func RespondWith404(w http.ResponseWriter, message string) error {
	return RespondWithError(w,
		http.StatusNotFound,
		message)
}

func RespondWith409(w http.ResponseWriter, message string) error {
	return RespondWithError(w,
		http.StatusConflict,
		message)
}

func RespondWith500(w http.ResponseWriter) error {
	return RespondWithError(w,
		http.StatusInternalServerError,
		http.StatusText(http.StatusInternalServerError))
}

func SuccessRespondWith200(w http.ResponseWriter, payload interface{}) error {
	return RespondWithJSON(w,
		http.StatusOK,
		payload)
}

func SuccessRespondWith201(w http.ResponseWriter, payload interface{}) error {
	return RespondWithJSON(w,
		http.StatusCreated,
		payload)
}
