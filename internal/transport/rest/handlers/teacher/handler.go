package teacher

import (
	"context"
	"net/http"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	teacherRoute  = "/teacher"
	teachersRoute = "/teachers"
)

type TeacherService interface {
	AddSkill(ctx context.Context, userID, categoryID int, videoCardLink string, about string) error
	BecomeTeacher(ctx context.Context, userID int) error
	GetTeacher(ctx context.Context, teacher *entities.Teacher) (*entities.User, error)
	GetTeacherList(ctx context.Context, userID int, isMyTeachers bool, category string, isFilteredByCategory bool) ([]entities.User, error)

	GetTeacherByID(ctx context.Context, id int) (*entities.Teacher, error)
	GetTeacherByUserID(ctx context.Context, userID int) (*entities.Teacher, error)
}

type TeacherHandlers struct {
	teacherService TeacherService
	log            *zap.Logger
}

func NewTeacherHandlers(teacherService TeacherService, log *zap.Logger) *TeacherHandlers {
	return &TeacherHandlers{
		teacherService: teacherService,
		log:            log,
	}
}

func (h *TeacherHandlers) SetupTeacherRoutes(router *chi.Mux, authMiddleware func(http.Handler) http.Handler) {
	teachersRouter := chi.NewRouter()

	teachersRouter.Get(getTeacherPublicRoute, h.GetTeacherPublic())

	teachersRouter.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		r.Get(getTeacherListRoute, h.GetTeacherList())

	})
	router.Mount(teachersRoute, teachersRouter)

	teacherRouter := chi.NewRouter()

	teacherRouter.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		r.Post(addSkillRoute, h.AddSkill())
		r.Post(becomeRoute, h.BecomeTeacher())
		r.Get(getTeacherProtectedRoute, h.GetTeacherProtected())
	})

	router.Mount(teacherRoute, teacherRouter)
}
