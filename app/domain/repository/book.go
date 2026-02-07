package repository

import (
	"context"
	"errors"

	"github.com/sora-00/booktracker-api/app/domain/entity"
)

// BookRepo は本の永続化のインターフェース。
// PostgreSQL などの実装は app/infra/repository/postgres にあります。
type BookRepo interface {
	Create(ctx context.Context, book *entity.Book) error
	FindAll(ctx context.Context) ([]entity.Book, error)
	FindByID(ctx context.Context, id int) (*entity.Book, error)
	Delete(ctx context.Context, id int) error
}

// ErrNotFound は対象が存在しないときに返す。controller で 404 に変換する。
var ErrNotFound = errors.New("not found")
