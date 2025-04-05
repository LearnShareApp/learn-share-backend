package teacher

import (
	"context"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

func (s *TeacherService) GetTeacherList(ctx context.Context, userID int, isMyTeachers bool, category string, isFilteredByCategory bool) ([]entities.User, error) {
	teachers, err := s.repo.GetAllTeachersDataFiltered(ctx, userID, isMyTeachers, category, isFilteredByCategory)
	if err != nil {
		return nil, fmt.Errorf("failed to get list of teachers: %w", err)
	}

	return teachers, nil
}
