package teacher

import (
	"context"
	"errors"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

func (s *TeacherService) GetTeacher(ctx context.Context, teacher *entities.Teacher) (*entities.User, error) {
	// get teacher's user data (common data)
	user, err := s.repo.GetUserByID(ctx, teacher.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher's user data: %w", err)
	}

	stat, err := s.repo.GetShortStatTeacherByID(ctx, teacher.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stat teacher: %w", err)
	}

	teacher.TeacherStat = *stat

	user.TeacherData = teacher

	teacher.Skills, err = s.repo.GetSkillsByTeacherID(ctx, teacher.ID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return user, nil
		}

		return nil, fmt.Errorf("failed to get skills: %w", err)
	}

	return user, nil
}
