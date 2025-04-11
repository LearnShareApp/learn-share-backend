package admin

import "context"

type Repository interface {
	IsUserAdminByID(ctx context.Context, id int) (bool, error)
}

type AdminService struct {
	repo Repository
}

func NewService(repo Repository) *AdminService {
	return &AdminService{
		repo: repo,
	}
}
