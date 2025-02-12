package get_teacher

import (
	"context"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	GetTeacherByUserId(ctx context.Context, id int) (*entities.Teacher, error)
	GetTeacherById(ctx context.Context, id int) (*entities.Teacher, error)
	GetUserById(ctx context.Context, id int) (*entities.User, error)
	GetShortStatTeacherById(ctx context.Context, id int) (*entities.TeacherStatistic, error)
	GetSkillsByTeacherId(ctx context.Context, id int) ([]*entities.Skill, error)
}
