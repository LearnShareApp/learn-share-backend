package admin

import (
	"context"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type Repository interface {
	IsUserAdminByID(ctx context.Context, id int) (bool, error)

	GetSkillByID(ctx context.Context, id int) (*entities.Skill, error)
	ActivateSkillByID(ctx context.Context, id int) error
}

type AdminService struct {
	repo Repository
}

func NewService(repo Repository) *AdminService {
	return &AdminService{
		repo: repo,
	}
}
