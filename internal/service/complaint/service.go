package complaint

import (
	"context"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type Repository interface {
	CreateComplaint(ctx context.Context, complaint *entities.Complaint) error
	GetAllComplaints(ctx context.Context) ([]*entities.Complaint, error)

	IsUserExistsByID(ctx context.Context, userId int) (bool, error)
	GetUserByID(ctx context.Context, id int) (*entities.User, error)
}

type ComplaintService struct {
	repo Repository
}

func NewService(repo Repository) *ComplaintService {
	return &ComplaintService{
		repo: repo,
	}
}
