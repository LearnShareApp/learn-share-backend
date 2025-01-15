package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const UserIDKey = "user_id"

var (
	ErrorTokenExpired = errors.New("token is expired")
)

type Service struct {
	secretKey []byte
	issuer    string
	duration  time.Duration
}

type Option func(*Service)

func WithDuration(duration time.Duration) Option {
	return func(s *Service) {
		s.duration = duration
	}
}

func WithIssuer(issuer string) Option {
	return func(s *Service) {
		s.issuer = issuer
	}
}

func NewJwtService(secretKey string, opts ...Option) *Service {
	s := &Service{
		secretKey: []byte(secretKey),
		duration:  24 * time.Hour,
		issuer:    "default",
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// GenerateJWTToken создает JWT-токен для пользователя
func (s *Service) GenerateJWTToken(userId int) (string, error) {
	// Устанавливаем время жизни токена
	expirationTime := time.Now().Add(s.duration)

	// Создаем claims
	claims := jwt.MapClaims{
		UserIDKey: userId,
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
		"iss":     s.issuer,
	}

	// Создаем токен с алгоритмом подписи
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключом
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

// ValidateJWTToken проверяет валидность JWT-токена
func (s *Service) ValidateJWTToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Дополнительная проверка времени истечения
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, ErrorTokenExpired
		}
	}

	return claims, nil
}

// ExtractUserID извлекает ID пользователя из claims
func (s *Service) ExtractUserID(claims jwt.MapClaims) (int, error) {
	userID, ok := claims[UserIDKey].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid or missing user ID in claims")
	}
	return int(userID), nil
}

func (s *Service) GetUserKey() string {
	return UserIDKey
}

func (s *Service) GetExpiredError() error {
	return ErrorTokenExpired
}
