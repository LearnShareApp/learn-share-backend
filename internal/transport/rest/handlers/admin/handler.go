package admin

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
)

const (
	adminRoute = "/admin"
)

type AdminService interface {
	CheckUserOnAdminByID(ctx context.Context, id int) (bool, error)
	ApproveTeacherSkill(ctx context.Context, skillID int) error
}

type AdminHandlers struct {
	service AdminService
	log     *zap.Logger
}

func NewAdminHandlers(adminService AdminService, log *zap.Logger) *AdminHandlers {
	return &AdminHandlers{
		service: adminService,
		log:     log,
	}
}

func (h *AdminHandlers) SetupAdminRoutes(router *chi.Mux, authMiddleware func(http.Handler) http.Handler) {
	adminRouter := chi.NewRouter()

	adminRouter.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		r.Get(checkRoute, h.CheckOnAdmin())
		r.Put(approveSkillRoute, h.ApproveSkill())
	})

	router.Mount(adminRoute, adminRouter)
}
