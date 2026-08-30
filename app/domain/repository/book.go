package repository

import (
	"context"

	"github.com/sora-00/booktracker-api/app/domain/entity"
)

type BookRepo interface {
	Put(ctx context.Context, book *entity.Book) error
	FindAll(ctx context.Context, userID string) ([]entity.Book, error)
	FindByStatus(ctx context.Context, userID string, status entity.Status) ([]entity.Book, error)
	FindByID(ctx context.Context, userID string, id int) (*entity.Book, error)
	Delete(ctx context.Context, userID string, id int) error
}
