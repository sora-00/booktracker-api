package repository

import (
	"context"

	"github.com/sora-00/booktracker-api/app/domain/entity"
)

type LogRepo interface {
	Put(ctx context.Context, log *entity.Log) error
	FindAll(ctx context.Context, userID string) ([]entity.Log, error)
	FindByID(ctx context.Context, userID string, id int) (*entity.Log, error)
	FindByBookID(ctx context.Context, userID string, bookID int) ([]entity.Log, error)
	DeleteByBookID(ctx context.Context, userID string, bookID int) error
	DeleteByLogID(ctx context.Context, userID string, id int) error
}
