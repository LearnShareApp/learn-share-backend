package repository

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"time"
)

type LearnShareRepo interface { //nolint:interfacebloat
	IsUserExistsById(ctx context.Context, id int) (bool, error)
	IsUserExistsByEmail(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, user *entities.User) (int, error)
	UpdateUser(ctx context.Context, userID int, user *entities.User) error
	GetUserById(ctx context.Context, id int) (*entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	GetUserStatByUserId(gs context.Context, id int) (*entities.StudentStatistic, error)

	IsTeacherExistsByUserId(ctx context.Context, id int) (bool, error)
	IsTeacherExistsById(ctx context.Context, id int) (bool, error)
	CreateTeacher(ctx context.Context, userID int) error
	CreateTeacherIfNotExists(ctx context.Context, userID int) (int, error)
	GetTeacherByUserId(ctx context.Context, id int) (*entities.Teacher, error)
	GetTeacherById(ctx context.Context, id int) (*entities.Teacher, error)
	GetShortStatTeacherById(ctx context.Context, id int) (*entities.TeacherStatistic, error)
	GetSkillsByTeacherId(ctx context.Context, id int) ([]*entities.Skill, error) // TODO: remake next
	GetAllTeachersDataFiltered(ctx context.Context, userId int, isUsersTeachers bool, category string, isFilteredByCategory bool) ([]entities.User, error)

	IsScheduleTimeExistsByTeacherIdAndDatetime(ctx context.Context, id int, datetime time.Time) (bool, error)
	CreateScheduleTime(ctx context.Context, teacherId int, datetime time.Time) error
	GetScheduleTimeById(ctx context.Context, id int) (*entities.ScheduleTime, error)
	GetScheduleTimesByTeacherId(ctx context.Context, id int) ([]*entities.ScheduleTime, error)

	GetSkillIdByTeacherIdAndCategoryId(ctx context.Context, teacherId int, categoryId int) (int, error)
	CreateReview(ctx context.Context, review *entities.Review) error
	GetReviewsByTeacherId(ctx context.Context, id int) ([]*entities.Review, error)

	IsCategoryExistsById(ctx context.Context, id int) (bool, error)
	GetCategories(ctx context.Context) ([]*entities.Category, error)

	IsFinishedLessonExistsByTeacherIdAndStudentIdAndCategoryId(ctx context.Context, teacherID int, studentId int, categoryId int) (bool, error)
	CreateUnconfirmedLesson(ctx context.Context, lesson *entities.Lesson) error
	GetLessonById(ctx context.Context, id int) (*entities.Lesson, error)
	GetStatusIdByStatusName(ctx context.Context, name string) (int, error)
	GetStudentLessonsByUserId(ctx context.Context, id int) ([]*entities.Lesson, error)
	GetTeacherLessonsByTeacherId(ctx context.Context, id int) ([]*entities.Lesson, error)
	ChangeLessonStatus(ctx context.Context, lessonID int, statusID int) error

	IsSkillExistsByTeacherIdAndCategoryId(ctx context.Context, teacherId int, categoryId int) (bool, error)
	CreateSkill(ctx context.Context, skill *entities.Skill) error
}

type Repository struct {
	db         *sqlx.DB
	sqlBuilder squirrel.StatementBuilderType
}

func New(db *sqlx.DB) *Repository {

	return &Repository{
		db:         db,
		sqlBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
