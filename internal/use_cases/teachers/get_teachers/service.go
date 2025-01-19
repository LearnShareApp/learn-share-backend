package get_teachers

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type Service struct {
	repo repo
}

func NewService(repo repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Do(ctx context.Context) ([]entities.User, error) {

	// maybe add filters

	teachers, err := s.repo.GetAllTeachersData(ctx)

	if err != nil {
		return nil, err
	}
	return teachers, nil
}
