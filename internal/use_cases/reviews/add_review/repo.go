package add_review

import (
	"context"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type repo interface {
	IsUserExistsById(ctx context.Context, id int) (bool, error)
	GetTeacherById(ctx context.Context, id int) (*entities.Teacher, error)
	IsCategoryExistsById(ctx context.Context, id int) (bool, error)
	GetSkillIdByTeacherIdAndCategoryId(ctx context.Context, teacherId int, categoryId int) (int, error)
	IsFinishedLessonExistsByTeacherIdAndStudentIdAndCategoryId(ctx context.Context, teacherId int, studentId int, categoryId int) (bool, error)
	CreateReview(ctx context.Context, review *entities.Review) error
}
