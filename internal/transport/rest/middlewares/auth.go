package middlewares

import (
	"context"
	"errors"
	"github.com/LearnShareApp/learn-share-backend/internal/httputils"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type TokenValidator interface {
	ValidateJWTToken(tokenString string) (jwt.MapClaims, error)
	ExtractUserID(claims jwt.MapClaims) (int, error)
	GetUserKey() string
	GetExpiredError() error
}

// JWTMiddleware middleware для проверки JWT токена
func JWTMiddleware(validator TokenValidator, log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем заголовок Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				if err := httputils.RespondWith401(w, "missed Authorization header (required)"); err != nil {
					log.Error("failed to write response", zap.Error(err))
				}
				return
			}

			// Проверяем формат заголовка (Bearer Token)
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				if err := httputils.RespondWith401(w, "Invalid token format"); err != nil {
					log.Error("failed to write response", zap.Error(err))
				}
				return
			}

			// Извлекаем токен
			tokenString := parts[1]

			// Валидируем токен с помощью переданного валидатора
			claims, err := validator.ValidateJWTToken(tokenString)
			if err != nil {
				if errors.Is(err, validator.GetExpiredError()) {
					log.Error("token expired", zap.Error(err))
					if err = httputils.RespondWith401(w, "token expired"); err != nil {
						log.Error("failed to write response", zap.Error(err))
					}
					return
				}

				log.Error("failed to validate token", zap.Error(err))
				if err = httputils.RespondWith401(w, "Failed to validate token"); err != nil {
					log.Error("failed to write response", zap.Error(err))
				}
				return
			}

			// Извлекаем ID из токена
			userID, err := validator.ExtractUserID(claims)
			if err != nil {
				log.Error("failed to extract user ID", zap.Error(err))
				if err = httputils.RespondWith401(w, "Invalid token: missing field: user_id"); err != nil {
					log.Error("failed to write response", zap.Error(err))
				}
				return
			}

			ctx := context.WithValue(r.Context(), validator.GetUserKey(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
