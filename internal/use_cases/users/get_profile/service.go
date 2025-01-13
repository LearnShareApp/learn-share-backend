package get_profile

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
)

type Service struct {
	repo repo
}

func NewService(repo repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Do(ctx context.Context, id int64) (*entities.User, error) {

	// Вроде как нет смысла обрабатывать кейс когда пользователь не найден по id и заворачивать в ошибку
	// т. к. по хорошему в токене, который выпустили мы не может быть несуществующий пользователь

	return s.repo.GetUserById(ctx, id)
}
