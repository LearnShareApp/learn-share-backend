package teacher

import (
	"context"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type Repository interface {
	IsUserExistsByID(ctx context.Context, id int) (bool, error)
	IsCategoryExistsByID(ctx context.Context, id int) (bool, error)
	CreateTeacherIfNotExists(ctx context.Context, userId int) (int, error)
	CreateSkill(ctx context.Context, skill *entities.Skill) error
	IsTeacherExistsByUserID(ctx context.Context, id int) (bool, error)
	CreateTeacher(ctx context.Context, userID int) error

	GetTeacherByUserID(ctx context.Context, id int) (*entities.Teacher, error)
	GetTeacherByID(ctx context.Context, id int) (*entities.Teacher, error)
	GetUserByID(ctx context.Context, id int) (*entities.User, error)
	GetShortStatTeacherByID(ctx context.Context, id int) (*entities.TeacherStatistic, error)
	GetSkillsByTeacherID(ctx context.Context, id int) ([]*entities.Skill, error)

	GetAllTeachersDataFiltered(ctx context.Context, userID int, isUsersTeachers bool, category string, isFilteredByCategory bool) ([]entities.User, error)
}

type TeacherService struct {
	repo Repository
}

func NewService(repo Repository) *TeacherService {
	return &TeacherService{
		repo: repo,
	}
}
