package schedule

import (
	"context"
	"net/http"
	"path"
	"time"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	teacherRoute  = "/teacher"
	teachersRoute = "/teachers"
)

type ScheduleService interface {
	AddTime(ctx context.Context, userID int, datetime time.Time) error
	GetTimes(ctx context.Context, teacher *entities.Teacher) ([]*entities.ScheduleTime, error)

	GetTeacherByID(ctx context.Context, id int) (*entities.Teacher, error)
	GetTeacherByUserID(ctx context.Context, userID int) (*entities.Teacher, error)
}

type ScheduleHandlers struct {
	scheduleService ScheduleService
	log             *zap.Logger
}

func NewScheduleHandlers(scheduleService ScheduleService, log *zap.Logger) *ScheduleHandlers {
	return &ScheduleHandlers{
		scheduleService: scheduleService,
		log:             log,
	}
}

func (h *ScheduleHandlers) SetupScheduleRoutes(router *chi.Mux, authMiddleware func(http.Handler) http.Handler) {
	router.Get(path.Join(teachersRoute, PublicGetListRoute), h.GetSchedulePublic())

	router.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post(path.Join(teacherRoute, AddRoute), h.AddScheduleTime())
		r.Get(path.Join(teacherRoute, ProtectedGetListRoute), h.GetScheduleProtected())
	})
}
