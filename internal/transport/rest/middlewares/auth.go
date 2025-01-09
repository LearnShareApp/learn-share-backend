package middlewares

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/jsonutils"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type TokenValidator interface {
	ValidateJWTToken(tokenString string) (jwt.MapClaims, error)
	ExtractUserID(claims jwt.MapClaims) (int64, error)
	GetUserKey() string
}

// JWTMiddleware middleware для проверки JWT токена
func JWTMiddleware(validator TokenValidator, log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Получаем заголовок Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				if err := jsonutils.RespondWith401(w, "missed Authorization header (required)"); err != nil {
					log.Error(err.Error())
					return
				}

				// Проверяем формат заголовка (Bearer Token)
				parts := strings.Split(authHeader, " ")
				if len(parts) != 2 || parts[0] != "Bearer" {
					if err := jsonutils.RespondWith401(w, "Invalid token format"); err != nil {
						log.Error(err.Error())
					}
					return
				}

				// Извлекаем токен
				tokenString := parts[1]

				// Валидируем токен с помощью переданного валидатора
				claims, err := validator.ValidateJWTToken(tokenString)

				if err != nil {
					log.Error(err.Error())
					if err = jsonutils.RespondWith401(w, "Failed to validate token"); err != nil {
						log.Error(err.Error())
					}
					return
				}

				// Извлекаем ID из токена
				userID, err := validator.ExtractUserID(claims)
				if err != nil {
					log.Info(err.Error())
					if err = jsonutils.RespondWith401(w, "Invalid token: missing filed: user_id"); err != nil {
						log.Error(err.Error())
					}
				}

				ctx := context.WithValue(r.Context(), validator.GetUserKey(), userID)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}
