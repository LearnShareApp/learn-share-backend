package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Список чувствительных полей, которые нужно маскировать
var sensitiveFields = []string{
	"password",
	"token",
	"jwt",
	"authorization",
	"api_key",
	"access_token",
	"refresh_token",
	"credit_card",
	"card_number",
}

type responseWriter struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		body:           bytes.NewBuffer([]byte{}),
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(body []byte) (int, error) {
	rw.body.Write(body)
	return rw.ResponseWriter.Write(body)
}

// maskSensitiveData маскирует чувствительные данные в JSON
func maskSensitiveData(data []byte) string {
	// Пробуем распарсить как JSON
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		// Если это не JSON, просто проверяем наличие чувствительных полей
		bodyStr := string(data)
		for _, field := range sensitiveFields {
			if strings.Contains(strings.ToLower(bodyStr), field) {
				return "[FILTERED]"
			}
		}
		return bodyStr
	}

	// Рекурсивно маскируем чувствительные поля
	maskJSONFields(jsonMap)

	// Преобразуем обратно в JSON
	maskedJSON, err := json.Marshal(jsonMap)
	if err != nil {
		return "[FILTERED]"
	}
	return string(maskedJSON)
}

// maskJSONFields рекурсивно маскирует чувствительные поля в JSON-структуре
func maskJSONFields(data map[string]interface{}) {
	for key, value := range data {
		lowerKey := strings.ToLower(key)

		// Проверяем, является ли поле чувствительным
		for _, sensitive := range sensitiveFields {
			if strings.Contains(lowerKey, sensitive) {
				data[key] = "[FILTERED]"
				break
			}
		}

		// Рекурсивно обрабатываем вложенные объекты
		switch v := value.(type) {
		case map[string]interface{}:
			maskJSONFields(v)
		case []interface{}:
			for _, item := range v {
				if mapItem, ok := item.(map[string]interface{}); ok {
					maskJSONFields(mapItem)
				}
			}
		}
	}
}

func LoggerMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Маскируем чувствительные заголовки
			headers := make(map[string]string)
			for key, values := range r.Header {
				lowerKey := strings.ToLower(key)
				if slices.Contains(sensitiveFields, lowerKey) {
					headers[key] = "[FILTERED]"
				} else {
					headers[key] = strings.Join(values, ", ")
				}
			}

			logger.Info("request started",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
			)

			// Читаем и маскируем тело запроса
			var requestBody []byte
			if r.Body != nil {
				requestBody, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
			}

			rw := newResponseWriter(w)
			next.ServeHTTP(rw, r)

			duration := time.Since(start)

			// Формируем лог с маскированными данными
			logFields := []zap.Field{
				zap.Duration("duration", duration),
				zap.Int("status", rw.status),
			}

			// логируем request body / response body только в случае ошибки
			if rw.status >= 400 {
				maskedRequest := maskSensitiveData(requestBody)
				logFields = append(logFields,
					zap.String("error_code", http.StatusText(rw.status)),
					zap.String("request_body", maskedRequest),
				)

				// Маскируем тело ответа
				if rw.body.Len() > 0 {
					maskedResponse := maskSensitiveData(rw.body.Bytes())
					logFields = append(logFields, zap.String("response_body", maskedResponse))
				}
			}

			logger.Info("request completed", logFields...)
		})
	}
}
