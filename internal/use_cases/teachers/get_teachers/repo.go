package get_teachers

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	GetAllTeachersDataFiltered(ctx context.Context, userId int, isUsersTeachers bool, category string, isFilteredByCategory bool) ([]entities.User, error)
}
