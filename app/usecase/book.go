package usecase

import (
	"context"
	"time"

	"github.com/sora-00/booktracker-api/app/domain/entity"
	"github.com/sora-00/booktracker-api/app/domain/repository"
	"github.com/sora-00/booktracker-api/app/domain/service"
	"github.com/sora-00/booktracker-api/app/usecase/request"
)

// BookUsecase は book に関するユースケースをまとめる
type BookUsecase struct {
	bookRepo    repository.BookRepo // DBアクセス
	bookService *service.BookSvc      // ドメインロジック
}

func NewUsecase(repo repository.BookRepo, svc *service.BookSvc) *BookUsecase {
	return &BookUsecase{
		bookRepo:    repo,
		bookService: svc,
	}
}

// すべての本を取得
func (u *BookUsecase) GetAllBooks(ctx context.Context) ([]entity.Book, error) {
	return u.bookRepo.FindAll(ctx)
}

// ID指定で1冊取得
func (u *BookUsecase) GetBookByID(ctx context.Context, id int) (*entity.Book, error) {
	return u.bookRepo.FindByID(ctx, id)
}

// 本を登録（request を entity に変換するのは usecase の役割）
func (u *BookUsecase) CreateBook(ctx context.Context, req *request.BookCreate) (*entity.Book, error) {
	now := time.Now()
	book := &entity.Book{
		Title:              req.Title,
		Author:             req.Author,
		TotalPages:         req.TotalPages,
		Publisher:          req.Publisher,
		ThumbnailUrl:       req.ThumbnailUrl,
		Status:             entity.Status(req.Status),
		TargetCompleteDate: req.TargetCompleteDate,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	return u.bookService.CreateBook(ctx, book)
}

// 本を削除
func (u *BookUsecase) DeleteBook(ctx context.Context, id int) error {
	return u.bookRepo.Delete(ctx, id)
}
