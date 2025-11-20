package service

import (
	"context"
	"errors"
	"time"

	"github.com/sora-00/booktracker-api/app/domain/entity"
	"github.com/sora-00/booktracker-api/app/domain/repository"
)

type BookSvc struct {
	repo repository.BookRepo // 抽象interfaceに依存
}

func NewService(repo repository.BookRepo) *BookSvc {
	return &BookSvc{repo: repo}
}

// CreateBook はタイトル重複を避けつつ新しい本を登録するビジネスルール
func (s *BookSvc) CreateBook(ctx context.Context, book *entity.Book) (*entity.Book, error) {
	if book == nil {
		return nil, errors.New("book is required")
	}
	if book.Title == "" {
		return nil, errors.New("title is required")
	}

	book.CreatedAt = time.Now()

	if err := s.repo.Create(ctx, book); err != nil {
		return nil, err
	}

	return book, nil
}
