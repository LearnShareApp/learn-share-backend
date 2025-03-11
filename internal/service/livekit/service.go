package livekit

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/livekit/protocol/auth"
)

// LiveKitConfig contains LiveKit API credentials.
type LiveKitConfig struct {
	ApiKey    string `env:"LIVEKIT_API_KEY"`
	ApiSecret string `env:"LIVEKIT_API_SECRET"`
}

// Service handles LiveKit operations.
type Service struct {
	ApiKey    string
	ApiSecret string
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
func NewService(config LiveKitConfig, opts ...Option) *Service {
	s := &Service{
		ApiKey:    config.ApiKey,
		ApiSecret: config.ApiSecret,
		duration:  24 * time.Hour,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// GenerateMeetingToken creates a new LiveKit access token for a room.
func (s *Service) GenerateMeetingToken(roomName string) (string, error) {
	canPublishMedia := true
	canSubscribeMedia := true
	// Common token for all room participants
	grants := &auth.VideoGrant{
		RoomJoin:     true,
		Room:         roomName,
		CanPublish:   &canPublishMedia,
		CanSubscribe: &canSubscribeMedia,
	}

	at := auth.NewAccessToken(s.ApiKey, s.ApiSecret)
	at.SetVideoGrant(grants)

	// ADD this - unique identifier
	at.SetIdentity(fmt.Sprintf("participant_%s", uuid.New().String()))
	// Can set token lifetime
	at.SetValidFor(s.duration)

	return at.ToJWT()
}

// NameRoomByLessonId generates a room name based on lesson ID.
func (s *Service) NameRoomByLessonId(lessonId int) string {
	return fmt.Sprintf("lesson_#%d", lessonId)
}
