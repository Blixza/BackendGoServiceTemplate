package town_service

import (
	town "backend-service-template/internal/domain/town"
	"context"

	"github.com/google/uuid"
)

type Service struct {
	repo town.Repository
}

func NewService(repo town.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(
	ctx context.Context, name string, balance int,
	ownerNickname string,
	xCoordOverworld int, yCoordOverworld int, zCoordOverworld int,
	xCoordNether int, yCoordNether int, zCoordNether int,
) (*town.Model, error) {
	town := &town.Model{
		Name:            name,
		Balance:         balance,
		OwnerNickname:   ownerNickname,
		XCoordOverworld: xCoordNether,
		YCoordOverworld: yCoordOverworld,
		ZCoordOverworld: zCoordOverworld,
		XCoordNether:    xCoordNether,
		YCoordNether:    yCoordNether,
		ZCoordNether:    zCoordNether,
	}

	err := s.repo.Create(ctx, town)
	if err != nil {
		return nil, err
	}

	return town, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*town.Model, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByName(ctx context.Context, name string) (*town.Model, error) {
	return s.repo.GetByName(ctx, name)
}

func (s *Service) GetByOwner(ctx context.Context, nickname string) (*town.Model, error) {
	return s.repo.GetByOwner(ctx, nickname)
}
