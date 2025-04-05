package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/LearnShareApp/learn-share-backend/pkg/hasher"
)

// CheckUser check user existence (by email) and compare password, if all correct returns his id.
func (s *UserService) CheckUser(ctx context.Context, reqUser *entities.User) (int, error) {
	realUser, err := s.repo.GetUserByEmail(ctx, reqUser.Email)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return 0, serviceErrs.ErrorUserNotFound
		}

		return 0, fmt.Errorf("failed to get user: %w", err)
	}

	if !hasher.ComparePassword(reqUser.Password, realUser.Password) {
		return 0, serviceErrs.ErrorPasswordIncorrect
	}

	return realUser.Id, nil
}
