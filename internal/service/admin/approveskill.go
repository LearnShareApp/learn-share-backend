package admin

import (
	"context"
	"errors"
	"fmt"

	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

func (s *AdminService) ApproveTeacherSkill(ctx context.Context, skillID int) error {
	skill, err := s.repo.GetSkillByID(ctx, skillID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorSkillNotFound
		}

		return fmt.Errorf("failed to get skill by id: %w", err)
	}

	if skill.IsActive {
		return serviceErrs.ErrorSkillAlreadyApproved
	}

	if err = s.repo.ActivateSkillByID(ctx, skillID); err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorSkillNotFound
		}

		return fmt.Errorf("failed to activate skill by id: %w", err)
	}

	return nil
}
