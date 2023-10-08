package repository

import (
	"context"

	"github.com/tahmooress/weConnect-task/internal/entity"
)

type Repository interface {
	Insert(ctx context.Context, record *entity.Statistics) (string, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*entity.Statistics, error)
	GetAll(ctx context.Context, page, limit int64) ([]entity.Statistics, error)
	Close() error
}
