package usecase

import (
	"github.com/sora-00/booktracker-api/app/domain/entity"
	"github.com/sora-00/booktracker-api/app/domain/repository"
	"github.com/sora-00/booktracker-api/app/domain/service"
)

// BookUsecase は book に関するユースケースをまとめる
type BookUsecase struct {
	bookRepo    repository.Book        // DBアクセス
	bookService *service.BookService   // ドメインロジック
}

func NewUsecase(repo repository.Book, service *service.BookService) *BookUsecase {
	return &BookUsecase{
		bookRepo:    repo,
		bookService: service,
	}
}

// すべての本を取得
func (u *BookUsecase) GetAllBooks() ([]entity.Book, error) {
	return u.bookRepo.FindAll()
}

// ID指定で1冊取得
func (u *BookUsecase) GetBookByID(id int) (*entity.Book, error) {
	return u.bookRepo.FindByID(id)
}

// 本を登録
func (u *BookUsecase) CreateBook(title, author string) (*entity.Book, error) {
	book, err := u.bookService.CreateBook(title, author)
	if err != nil {
		return nil, err
	}
	return book, nil
}

// 本を削除
func (u *BookUsecase) DeleteBook(id int) error {
	return u.bookRepo.Delete(id)
}
