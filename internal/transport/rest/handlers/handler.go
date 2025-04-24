package handlers

import (
	"net/http"

	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/handlers/admin"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/handlers/complaint"

	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/handlers/category"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/handlers/image"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/handlers/lesson"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/handlers/review"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/handlers/schedule"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/handlers/teacher"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/handlers/user"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Services interface {
	user.JwtService
	user.UserService
	teacher.TeacherService
	schedule.ScheduleService
	review.ReviewService
	lesson.LessonService
	image.ImageService
	category.CategoryService
	complaint.ComplaintService
	admin.AdminService
}

type Handlers struct {
	services Services
	log      *zap.Logger
}

func NewHandlers(services Services, log *zap.Logger) *Handlers {
	return &Handlers{
		services: services,
		log:      log,
	}
}

func (h *Handlers) SetupRoutes(router *chi.Mux, authMiddleware func(http.Handler) http.Handler) {
	//recomendation from AI about downcast to certain interfaces (ISP)

	var userService user.UserService = h.services
	var jwtService user.JwtService = h.services
	userHandlers := user.NewUserHandlers(userService, jwtService, h.log)
	userHandlers.SetupUserRoutes(router, authMiddleware)

	var teacherService teacher.TeacherService = h.services
	teacherHandlers := teacher.NewTeacherHandlers(teacherService, h.log)
	teacherHandlers.SetupTeacherRoutes(router, authMiddleware)

	var scheduleService schedule.ScheduleService = h.services
	scheduleHandlers := schedule.NewScheduleHandlers(scheduleService, h.log)
	scheduleHandlers.SetupScheduleRoutes(router, authMiddleware)

	var reviewService review.ReviewService = h.services
	reviewHandlers := review.NewReviewHandlers(reviewService, h.log)
	reviewHandlers.SetupReviewRoutes(router, authMiddleware)

	var lessonService lesson.LessonService = h.services
	lessonHandlers := lesson.NewLessonHandlers(lessonService, h.log)
	lessonHandlers.SetupLessonRoutes(router, authMiddleware)

	var imageService image.ImageService = h.services
	imageHandlers := image.NewImageHandlers(imageService, h.log)
	imageHandlers.SetupImageRoutes(router)

	var categoryService category.CategoryService = h.services
	categoryHandlers := category.NewCategoryHandlers(categoryService, h.log)
	categoryHandlers.SetupCategoryRoutes(router)

	var complaintService complaint.ComplaintService = h.services
	complaintHandlers := complaint.NewComplaintHandlers(complaintService, h.log)
	complaintHandlers.SetupComplaintRoutes(router, authMiddleware)

	var adminService admin.AdminService = h.services
	adminHandlers := admin.NewAdminHandlers(adminService, h.log)
	adminHandlers.SetupAdminRoutes(router, authMiddleware)
}
