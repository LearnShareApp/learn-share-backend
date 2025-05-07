package teacher

import (
	"context"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

func (s *TeacherService) GetTeacherShortDataListByIDs(ctx context.Context, TeacherIDs []int) ([]entities.User, error) {
	// for reasure about unique
	ids := make(map[int]bool, len(TeacherIDs))
	for _, v := range TeacherIDs {
		ids[v] = true
	}

	teachers, err := s.repo.GetShortTeacherDatasByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get sort data list of teachers by IDs: %w", err)
	}

	return teachers, nil
}
