package user

import (
	"context"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	"github.com/LearnShareApp/learn-share-backend/pkg/storage/object"
)

type ObjectStorage interface {
	UploadFile(ctx context.Context, file *object.File) error
}

type Repository interface {
	IsUserExistsByEmail(ctx context.Context, email string) (bool, error)
	IsTeacherExistsByUserID(ctx context.Context, id int) (bool, error)
	GetUserByID(ctx context.Context, id int) (*entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	GetUserStatByUserID(ctx context.Context, id int) (*entities.StudentStatistic, error)
	UpdateUser(ctx context.Context, userID int, user *entities.User) error
	CreateUser(ctx context.Context, user *entities.User) (int, error)
}

type UserService struct {
	repo          Repository
	objectStorage ObjectStorage
}

func NewService(repo Repository, objectStorage ObjectStorage) *UserService {
	return &UserService{
		repo:          repo,
		objectStorage: objectStorage,
	}
}
