package admin

import (
	"context"
	"errors"
	"fmt"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

func (s *AdminService) CheckUserOnAdminByID(ctx context.Context, id int) (bool, error) {
	isAdmin, err := s.repo.IsUserAdminByID(ctx, id)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return false, serviceErrs.ErrorUserNotFound
		}

		return false, fmt.Errorf("failed to check is user an admin: %w", err)
	}

	return isAdmin, nil
}
