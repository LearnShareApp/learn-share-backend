package lesson

import (
	"context"
	"net/http"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	lessonsRoute = "/lessons"
)

type LessonService interface {
	BookLesson(ctx context.Context, lesson *entities.Lesson) error
	PlanLesson(ctx context.Context, userID int, lessonID int) error
	RejectLesson(ctx context.Context, userID int, lessonID int) error
	CancelLesson(ctx context.Context, userID int, lessonID int) error
	FinishLesson(ctx context.Context, userID int, lessonID int) error
	JoinLesson(ctx context.Context, userID int, lessonID int) (string, error)
	StartLesson(ctx context.Context, userID, lessonID int) (string, error)
	GetLesson(ctx context.Context, lessonID int) (*entities.Lesson, error)
	GetLessonShortData(ctx context.Context, lessonID int) (*entities.Lesson, error)
	GetStudentLessonList(ctx context.Context, userID int) ([]*entities.Lesson, error)
	GetTeacherLessonList(ctx context.Context, userID int) ([]*entities.Lesson, error)

	GetTeacherByID(ctx context.Context, id int) (*entities.Teacher, error)
	GetTeacherByUserID(ctx context.Context, userID int) (*entities.Teacher, error)
}

type LessonHandlers struct {
	lessonService LessonService
	log           *zap.Logger
}

func NewLessonHandlers(lessonService LessonService, log *zap.Logger) *LessonHandlers {
	return &LessonHandlers{
		lessonService: lessonService,
		log:           log,
	}
}

func (h *LessonHandlers) SetupLessonRoutes(router *chi.Mux, authMiddleware func(http.Handler) http.Handler) {
	lessonsRouter := chi.NewRouter()

	lessonsRouter.Get(getRoute, h.GetLesson())
	lessonsRouter.Get(getShortRoute, h.GetLessonShortData())

	lessonsRouter.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		r.Put(planRoute, h.PlanLesson())
		r.Put(rejectRoute, h.RejectLesson())
		r.Put(cancelRoute, h.CancelLesson())
		r.Put(startRoute, h.StartLesson())
		r.Get(joinRoute, h.JoinToLesson())
		r.Put(finishRoute, h.FinishLesson())

	})

	router.Mount(lessonsRoute, lessonsRouter)

	router.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		r.Post(bookRoute, h.BookLesson())
		r.Get(getForStudentListRoute, h.GetForStudentList())
		r.Get(getForTeacherListRoute, h.GetForTeacherList())
	})

}
