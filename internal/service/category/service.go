package category

import (
	"context"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type Repository interface {
	GetCategories(ctx context.Context) ([]*entities.Category, error)
}

type CategoryService struct {
	repo Repository
}

func NewService(repo Repository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

// GetCategories returns all categories.
func (s *CategoryService) GetCategories(ctx context.Context) ([]*entities.Category, error) {
	categories, err := s.repo.GetCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories from db: %w", err)
	}

	return categories, nil
}
