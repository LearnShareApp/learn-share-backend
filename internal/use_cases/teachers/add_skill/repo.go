package add_skill

import (
	"context"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	IsUserExistsById(ctx context.Context, id int) (bool, error)
	IsCategoryExistsById(ctx context.Context, id int) (bool, error)
	CreateTeacherIfNotExists(ctx context.Context, userId int) (int, error)
	CreateSkill(ctx context.Context, skill *entities.Skill) error
}
