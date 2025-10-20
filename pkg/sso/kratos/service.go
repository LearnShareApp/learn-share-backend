package kratos

import (
	"context"
	"errors"
	"time"

	"github.com/LearnShareApp/learn-share-backend/pkg/sso"
	kratosClient "github.com/ory/kratos-client-go"
)

var (
	ErrSessionNotFound = errors.New("session cookie not found")
	ErrSessionInvalid  = errors.New("session is invalid")
	ErrSessionInactive = errors.New("session is not active")
)

type Config struct {
	AdminAPIURL    string        `env:"KRATOS_ADMIN_API_URL"`
	RequestTimeout time.Duration `env:"KRATOS_REQUEST_TIMEOUT" envDefault:"5s"`
}

type KratosService struct {
	client *kratosClient.APIClient
	config Config
}

func New(config Config) *KratosService {
	kratosConfig := kratosClient.NewConfiguration()
	kratosConfig.Servers = []kratosClient.ServerConfiguration{
		{
			URL: config.AdminAPIURL, // http://127.0.0.1:4434
		},
	}

	return &KratosService{
		client: kratosClient.NewAPIClient(kratosConfig),
		config: config,
	}
}

func (s *KratosService) GetSession(ctx context.Context, sessionCookie string) (*sso.SessionData, error) {
	if sessionCookie == "" {
		return nil, ErrSessionNotFound
	}

	ctx, cancel := context.WithTimeout(ctx, s.config.RequestTimeout)
	defer cancel()

	session, resp, err := s.client.FrontendAPI.ToSession(ctx).
		Cookie(sessionCookie).
		Execute()

	if err != nil {

		if resp != nil && resp.StatusCode == 401 {
			return nil, ErrSessionInactive
		}
		return nil, ErrSessionInvalid
	}

	if session == nil || session.Identity == nil {
		return nil, ErrSessionInvalid
	}

	if !*session.Active {
		return nil, ErrSessionInactive
	}

	sessionData := &sso.SessionData{
		IdentityID:    session.Identity.Id,
		IsActive:      *session.Active,
		EmailVerified: s.isEmailVerified(session.Identity),
		Email:         s.extractEmail(session.Identity),
	}

	return sessionData, nil
}

func (s *KratosService) ValidateSession(ctx context.Context, sessionCookie string) error {
	_, err := s.GetSession(ctx, sessionCookie)
	return err
}

// isEmailVerified проверяет подтверждение email из identity
func (s *KratosService) isEmailVerified(identity *kratosClient.Identity) bool {
	if identity.VerifiableAddresses == nil {
		return false
	}

	for _, addr := range identity.VerifiableAddresses {
		if addr.Verified && addr.Value != "" {
			return true
		}
	}
	return false
}

func (s *KratosService) extractEmail(identity *kratosClient.Identity) string {
	// first check verified addresses
	if identity.VerifiableAddresses != nil {
		for _, addr := range identity.VerifiableAddresses {
			if addr.Verified && addr.Value != "" {
				return addr.Value
			}
		}
	}

	// fallback on traits if verified addr not found
	if traits, ok := identity.Traits.(map[string]interface{}); ok {
		if email, ok := traits["email"].(string); ok {
			return email
		}
	}

	return ""
}
