package usecase

import (
	"context"
	"time"

	"github.com/sora-00/booktracker-api/app/domain/entity"
	"github.com/sora-00/booktracker-api/app/domain/repository"
	"github.com/sora-00/booktracker-api/app/domain/service"
	"github.com/sora-00/booktracker-api/app/usecase/request"
	"github.com/sora-00/booktracker-api/app/usecase/response"
)

type Book struct {
	bookRepo    repository.BookRepo
	bookService *service.BookSvc
}

func NewBook(repo repository.BookRepo, svc *service.BookSvc) *Book {
	return &Book{
		bookRepo:    repo,
		bookService: svc,
	}
}

func (b Book) Get(ctx context.Context, r *request.BookGet) (*response.BookGet, error) {
	books, err := b.bookRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return response.NewBookGet(books), nil
}

func (b Book) GetByID(ctx context.Context, r *request.BookGetByID) (*response.BookGetByID, error) {
	book, err := b.bookRepo.FindByID(ctx, r.BookID)
	if err != nil {
		return nil, err
	}
	return response.NewBookGetByID(book), nil
}

func (b Book) Create(ctx context.Context, r *request.BookCreate) (*response.BookCreate, error) {
	now := time.Now()
	book := &entity.Book{
		Title:              r.Title,
		Author:             r.Author,
		TotalPages:         r.TotalPages,
		Publisher:          r.Publisher,
		ThumbnailUrl:       r.ThumbnailUrl,
		Status:             entity.Status(r.Status),
		TargetCompleteDate: r.TargetCompleteDate,
		EncounterNote:      r.EncounterNote,
		ReadPages:          r.ReadPages,
		TargetPagesPerDay:  r.TargetPagesPerDay,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	created, err := b.bookService.CreateBook(ctx, book)
	if err != nil {
		return nil, err
	}
	return response.NewBookCreate(created), nil
}

func (b Book) Delete(ctx context.Context, r *request.BookDelete) (*response.BookDelete, error) {
	if err := b.bookRepo.Delete(ctx, r.BookID); err != nil {
		return nil, err
	}
	return response.NewBookDelete(r.BookID), nil
}
