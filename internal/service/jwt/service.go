package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	UserIDKey  = "user_id"
	defaultTTL = time.Hour * 24 * 7
)

var ErrorTokenExpired = errors.New("token is expired")

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

func NewService(secretKey string, opts ...Option) *Service {
	s := &Service{
		secretKey: []byte(secretKey),
		duration:  defaultTTL,
		issuer:    "default",
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// GenerateJWTToken creates a JWT token for a user.
func (s *Service) GenerateJWTToken(userId int) (string, error) {
	// Set token expiration time
	expirationTime := time.Now().Add(s.duration)

	// Create claims
	claims := jwt.MapClaims{
		UserIDKey: userId,
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
		"iss":     s.issuer,
	}

	// Create token with signing algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with secret key
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

// ValidateJWTToken validates the JWT token.
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
		return nil, errors.New("invalid token")
	}

	// Additional expiration time check
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, ErrorTokenExpired
		}
	}

	return claims, nil
}

// ExtractUserID extracts user ID from claims.
func (s *Service) ExtractUserID(claims jwt.MapClaims) (int, error) {
	userID, ok := claims[UserIDKey].(float64)
	if !ok {
		return 0, errors.New("invalid or missing user ID in claims")
	}

	return int(userID), nil
}

func (s *Service) GetUserKey() string {
	return UserIDKey
}

func (s *Service) GetExpiredError() error {
	return ErrorTokenExpired
}
