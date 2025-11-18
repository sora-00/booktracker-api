package service

import (
	"errors"

	"github.com/sora-00/booktracker-api/app/domain/entity"
	"github.com/sora-00/booktracker-api/app/domain/repository"
)

type BookService struct {
	repo repository.Book // 抽象interfaceに依存
}

func NewService(repo repository.Book) *BookService {
	return &BookService{repo: repo}
}

// CreateBook はタイトル重複を避けつつ新しい本を登録するビジネスルール
func (s *BookService) CreateBook(title, author string, totalPages int, publisher string) (*entity.Book, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	// タイトル重複チェック
	books, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	for _, b := range books {
		if b.Title == title {
			return nil, errors.New("book title already exists")
		}
	}

	// 重複がなければ新しい本を作成
	book := &entity.Book{
		Title:      title,
		Author:     author,
		TotalPages: totalPages,
		Publisher: publisher,
	}

	if err := s.repo.Save(book); err != nil {
		return nil, err
	}

	return book, nil
}
