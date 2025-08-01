package repository

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

//
//type LearnShareRepo interface { //nolint:interfacebloat
//	IsUserExistsByID(ctx context.Context, id int) (bool, error)
//	IsUserExistsByEmail(ctx context.Context, email string) (bool, error)
//	CreateUser(ctx context.Context, user *entities.User) (int, error)
//	UpdateUser(ctx context.Context, userID int, user *entities.User) error
//	GetUserByID(ctx context.Context, id int) (*entities.User, error)
//	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
//	GetUserStatByUserID(gs context.Context, id int) (*entities.StudentStatistic, error)
//
//	IsTeacherExistsByUserID(ctx context.Context, id int) (bool, error)
//	IsTeacherExistsById(ctx context.Context, id int) (bool, error)
//	CreateTeacher(ctx context.Context, userID int) error
//	CreateTeacherIfNotExists(ctx context.Context, userID int) (int, error)
//	GetTeacherByUserID(ctx context.Context, id int) (*entities.Teacher, error)
//	GetTeacherByID(ctx context.Context, id int) (*entities.Teacher, error)
//	GetShortStatTeacherByID(ctx context.Context, id int) (*entities.TeacherStatistic, error)
//	GetSkillsByTeacherID(ctx context.Context, id int) ([]*entities.Skill, error) // TODO: remake next
//	GetAllTeachersDataFiltered(ctx context.Context, userId int, isUsersTeachers bool, category string, isFilteredByCategory bool) ([]entities.User, error)
//
//	IsScheduleTimeExistsByTeacherIDAndDatetime(ctx context.Context, id int, datetime time.Time) (bool, error)
//	CreateScheduleTime(ctx context.Context, teacherId int, datetime time.Time) error
//	GetScheduleTimeByID(ctx context.Context, id int) (*entities.ScheduleTime, error)
//	GetScheduleTimesByTeacherID(ctx context.Context, id int) ([]*entities.ScheduleTime, error)
//
//	GetSkillIdByTeacherIdAndCategoryId(ctx context.Context, teacherId int, categoryId int) (int, error)
//	CreateReview(ctx context.Context, review *entities.Review) error
//	GetReviewsByTeacherId(ctx context.Context, id int) ([]*entities.Review, error)
//
//	IsCategoryExistsByID(ctx context.Context, id int) (bool, error)
//	GetCategories(ctx context.Context) ([]*entities.Category, error)
//
//	IsFinishedLessonExistsByTeacherIdAndStudentIdAndCategoryId(ctx context.Context, teacherID int, studentId int, categoryId int) (bool, error)
//	CreateUnconfirmedLesson(ctx context.Context, lesson *entities.Lesson) error
//	GetLessonByID(ctx context.Context, id int) (*entities.Lesson, error)
//	GetStatusIDByStatusName(ctx context.Context, name string) (int, error)
//	GetStudentLessonsByUserID(ctx context.Context, id int) ([]*entities.Lesson, error)
//	GetTeacherLessonsByTeacherID(ctx context.Context, id int) ([]*entities.Lesson, error)
//	ChangeLessonStatus(ctx context.Context, lessonID int, statusID int) error
//
//	IsSkillExistsByTeacherIDAndCategoryID(ctx context.Context, teacherId int, categoryId int) (bool, error)
//	CreateSkill(ctx context.Context, skill *entities.Skill) error
//}

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
