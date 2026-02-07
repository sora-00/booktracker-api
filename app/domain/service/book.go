package service

import (
	"context"
	"errors"

	"github.com/sora-00/booktracker-api/app/domain/entity"
	"github.com/sora-00/booktracker-api/app/domain/repository"
)

type BookSvc struct {
	repo repository.BookRepo // 抽象interfaceに依存
}

func NewService(repo repository.BookRepo) *BookSvc {
	return &BookSvc{repo: repo}
}

// CreateBook は新しい本を登録する（入力は request 層で検証済み）
func (s *BookSvc) CreateBook(ctx context.Context, book *entity.Book) (*entity.Book, error) {
	if book == nil {
		return nil, errors.New("book is required")
	}
	if err := s.repo.Create(ctx, book); err != nil {
		return nil, err
	}
	return book, nil
}
