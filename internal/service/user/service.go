package user_service

import (
	user "backend-service-template/internal/domain/user"
	"context"
	"time"
)

type Service struct {
	repo user.Repository
}

func NewService(repo user.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, name, email string) (*user.Model, error) {
	user := &user.Model{
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}

	err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Get(ctx context.Context, id int) (*user.Model, error) {
	return s.repo.GetByID(ctx, id)
}
