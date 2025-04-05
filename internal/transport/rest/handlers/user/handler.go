package user

import (
	"context"
	"io"
	"net/http"
	"path"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	authRoute  = "/auth"
	userRoute  = "/user"
	usersRoute = "/users"
)

type UserService interface {
	CreateUser(ctx context.Context, user *entities.User, avatarReader io.Reader, avatarSize int64) (int, error)
	GetUser(ctx context.Context, id int) (*entities.User, error)
	EditUser(ctx context.Context, userID int, user *entities.User, avatarReader io.Reader, avatarSize int64) error
	CheckUser(ctx context.Context, reqUser *entities.User) (int, error)
}

type JwtService interface {
	GenerateJWTToken(userID int) (string, error)
}

type UserHandlers struct {
	userService UserService
	jwtService  JwtService
	log         *zap.Logger
}

func NewUserHandlers(userService UserService, jwtService JwtService, log *zap.Logger) *UserHandlers {
	return &UserHandlers{
		userService: userService,
		jwtService:  jwtService,
		log:         log,
	}
}

func (h *UserHandlers) SetupUserRoutes(router *chi.Mux, authMiddleware func(http.Handler) http.Handler) {
	authRouter := chi.NewRouter()
	authRouter.Post(RegistrationRoute, h.RegistrationUser())
	authRouter.Post(LoginRoute, h.LoginUser())
	router.Mount(authRoute, authRouter)

	router.Get(path.Join(usersRoute, GetPublicRoute), h.GetUserPublic())

	router.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		r.Get(path.Join(userRoute, GetProtectedRoute), h.GetUserProtected())
		r.Patch(path.Join(userRoute, EditRoute), h.EditUser())

	})
}
