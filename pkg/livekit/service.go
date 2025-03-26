package livekit

import (
	"fmt"
	"time"

	"github.com/livekit/protocol/auth"
)

const (
	defaultTokenLifetime = 24 * time.Hour
)

// Config contains LiveKit API credentials.
type Config struct {
	APIKey    string `env:"LIVEKIT_API_KEY"`
	APISecret string `env:"LIVEKIT_API_SECRET"`
}

// Service handles LiveKit operations.
type Service struct {
	APIKey    string
	APISecret string
	duration  time.Duration
}

// Option is a function type to configure Service.
type Option func(*Service)

// WithDuration sets token duration for the service.
func WithDuration(duration time.Duration) Option {
	return func(s *Service) {
		s.duration = duration
	}
}

// NewService creates a new LiveKit service instance.
func NewService(config Config, opts ...Option) *Service {
	service := &Service{
		APIKey:    config.APIKey,
		APISecret: config.APISecret,
		duration:  defaultTokenLifetime,
	}

	for _, opt := range opts {
		opt(service)
	}

	return service
}

// GenerateMeetingToken creates a new LiveKit access token for a room.
func (s *Service) GenerateMeetingToken(roomName, userName string) (string, error) {
	canPublishMedia := true
	canSubscribeMedia := true
	// Common token for all room participants
	grants := &auth.VideoGrant{ //nolint:exhaustruct
		RoomJoin:     true,
		Room:         roomName,
		CanPublish:   &canPublishMedia,
		CanSubscribe: &canSubscribeMedia,
	}

	accessToken := auth.NewAccessToken(s.APIKey, s.APISecret)
	accessToken.SetVideoGrant(grants)

	// ADD this - unique identifier (name)
	accessToken.SetIdentity(userName)
	// Can set token lifetime
	accessToken.SetValidFor(s.duration)

	return accessToken.ToJWT() //nolint:wrapcheck
}

// NameRoomByLessonID generates a room name based on lesson ID.
func (s *Service) NameRoomByLessonID(lessonId int) string {
	return fmt.Sprintf("lesson_#%d", lessonId)
}

func (s *Service) GetUserIdentityString(userName, userSurname string, id int) string {
	return fmt.Sprintf("%s %s (%d)", userName, userSurname, id)
}
