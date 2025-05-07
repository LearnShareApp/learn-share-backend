package skill

import (
	"context"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

func (s *SkillService) GetUnactiveSkillList(ctx context.Context) ([]entities.Skill, error) {
	skills, err := s.repo.GetUnactiveSkills(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get list of skill: %w", err)
	}

	return skills, nil
}
