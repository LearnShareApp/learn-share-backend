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

func (s *Service) Do(ctx context.Context, userId int, isMyTeachers bool, category string, isFilteredByCategory bool) ([]entities.User, error) {

	teachers, err := s.repo.GetAllTeachersDataFiltered(ctx, userId, isMyTeachers, category, isFilteredByCategory)

	if err != nil {
		return nil, err
	}
	return teachers, nil
}
