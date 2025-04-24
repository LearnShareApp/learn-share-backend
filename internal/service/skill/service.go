package skill

import (
	"context"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type Repository interface {
	GetSkillByID(ctx context.Context, id int) (*entities.Skill, error)
	ActivateSkillByID(ctx context.Context, id int) error

	IsUserExistsByID(ctx context.Context, id int) (bool, error)
	IsCategoryExistsByID(ctx context.Context, id int) (bool, error)
	CreateTeacherIfNotExists(ctx context.Context, userId int) (int, error)
	CreateSkill(ctx context.Context, skill *entities.Skill) error
}

type SkillService struct {
	repo Repository
}

func NewService(repo Repository) *SkillService {
	return &SkillService{
		repo: repo,
	}
}
