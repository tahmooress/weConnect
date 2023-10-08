package service

import (
	"context"

	"github.com/tahmooress/weConnect-task/internal/entity"
	"github.com/tahmooress/weConnect-task/internal/repository"
)

type UseCase interface {
	Create(ctx context.Context, record *entity.Statistics) (string, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*entity.Statistics, error)
	GetAll(ctx context.Context, page, limit int64) ([]entity.Statistics, error)
}

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, record *entity.Statistics) (string, error) {
	return s.repo.Insert(ctx, record)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) GetByID(ctx context.Context, id string) (*entity.Statistics, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetAll(ctx context.Context, page, limit int64) ([]entity.Statistics, error) {
	return s.repo.GetAll(ctx, page, limit)
}
