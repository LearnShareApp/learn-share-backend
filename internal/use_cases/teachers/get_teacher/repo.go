package get_teacher

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	GetUserById(ctx context.Context, id int) (*entities.User, error)
	GetTeacherByUserId(ctx context.Context, id int) (*entities.Teacher, error)
	GetSkillsByTeacherId(ctx context.Context, id int) ([]*entities.Skill, error)
}
