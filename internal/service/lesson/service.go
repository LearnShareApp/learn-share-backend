package lesson

import (
	"context"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type Repository interface {
	IsUserExistsByID(ctx context.Context, id int) (bool, error)
	GetTeacherByUserID(ctx context.Context, id int) (*entities.Teacher, error)
	GetLessonByID(ctx context.Context, id int) (*entities.Lesson, error)
	GetStatusIDByStatusName(ctx context.Context, name string) (int, error)
	ChangeLessonStatus(ctx context.Context, lessonID int, statusID int) error

	GetTeacherByID(ctx context.Context, id int) (*entities.Teacher, error)
	IsCategoryExistsByID(ctx context.Context, id int) (bool, error)
	IsSkillExistsByTeacherIDAndCategoryID(ctx context.Context, teacherID int, categoryID int) (bool, error)
	GetScheduleTimeByID(ctx context.Context, id int) (*entities.ScheduleTime, error)
	CreateUnconfirmedLesson(ctx context.Context, lesson *entities.Lesson) error

	GetUserByID(ctx context.Context, id int) (*entities.User, error)

	GetStudentLessonsByUserID(ctx context.Context, id int) ([]*entities.Lesson, error)
	GetTeacherLessonsByTeacherID(ctx context.Context, id int) ([]*entities.Lesson, error)
}

type MeetCreator interface {
	GenerateMeetingToken(roomName string, userName string) (string, error)
	NameRoomByLessonID(lessonID int) string
	GetUserIdentityString(userName, userSurname string, id int) string
}

type LessonService struct {
	repo        Repository
	meetCreator MeetCreator
}

func NewService(repo Repository, meet MeetCreator) *LessonService {
	return &LessonService{
		repo:        repo,
		meetCreator: meet,
	}
}
