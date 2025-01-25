package livekit

import (
	"github.com/livekit/protocol/auth"
	"time"
)

type ApiConfig struct {
	ApiKey    string
	ApiSecret string
}

type Service struct {
	ApiKey    string
	ApiSecret string
	duration  time.Duration
}

type Option func(*Service)

func WithDuration(duration time.Duration) Option {
	return func(s *Service) {
		s.duration = duration
	}
}

func NewService(config ApiConfig, opts ...Option) *Service {
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

func (s *Service) GenerateMeetingToken(roomName string) (string, error) {
	canPublishMedia := true
	CanSubscribeMedia := true
	// Общий токен для всех участников комнаты
	grants := &auth.VideoGrant{
		RoomJoin:     true,
		Room:         roomName,
		CanPublish:   &canPublishMedia,
		CanSubscribe: &CanSubscribeMedia,
	}

	at := auth.NewAccessToken(s.ApiKey, s.ApiSecret)
	at.SetVideoGrant(grants)
	//at.AddGrant(grants)

	// Можно добавить время жизни токена
	at.SetValidFor(s.duration)

	return at.ToJWT()
}
