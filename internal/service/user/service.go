package user_service

import (
	user "backend-service-template/internal/domain/user"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo user.Repository
}

func NewService(repo user.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(
	ctx context.Context, nickname, password, discord, email string,
	balance int, towns sql.NullString,
) (*user.Model, error) {
	user := &user.Model{
		Nickname:  nickname,
		Email:     email,
		Discord:   discord,
		Password:  password,
		Balance:   balance,
		Towns:     towns,
		CreatedAt: time.Now(),
	}

	err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetByNickname(ctx context.Context, nickname string) (*user.Model, error) {
	return s.repo.GetByNickname(ctx, nickname)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*user.Model, error) {
	return s.repo.GetByID(ctx, id)
}
