package get_teacher

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	IsUserExistsById(ctx context.Context, id int) (bool, error)
	IsTeacherExistsByUserId(ctx context.Context, id int) (bool, error)
	GetUserById(ctx context.Context, id int) (*entities.User, error)
	GetTeacherByUserId(ctx context.Context, id int) (*entities.Teacher, error)
	GetShortStatTeacherById(ctx context.Context, id int) (*entities.TeacherStatistic, error)
	GetSkillsByTeacherId(ctx context.Context, id int) ([]*entities.Skill, error)
}
