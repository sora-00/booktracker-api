package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/sora-00/booktracker-api/app/domain/entity"
	"github.com/sora-00/booktracker-api/app/domain/repository"
	"github.com/sora-00/booktracker-api/app/infra/auth"
	"github.com/sora-00/booktracker-api/app/usecase/request"
	"github.com/sora-00/booktracker-api/app/usecase/response"
)

type Book struct {
	bookRepo repository.BookRepo
	logRepo  repository.LogRepo
}

func NewBook(bookRepo repository.BookRepo, logRepo repository.LogRepo) *Book {
	return &Book{bookRepo: bookRepo, logRepo: logRepo}
}

func (b Book) Get(ctx context.Context, _ *request.BookGet) (*response.BookGet, error) {
	userID := auth.UserIDFromContext(ctx)
	books, err := b.bookRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	return response.NewBookGet(books), nil
}

func (b Book) GetByStatus(ctx context.Context, r *request.BookGetByStatus) (*response.BookGet, error) {
	userID := auth.UserIDFromContext(ctx)
	books, err := b.bookRepo.FindByStatus(ctx, userID, entity.Status(r.Status))
	if err != nil {
		return nil, err
	}
	return response.NewBookGet(books), nil
}

func (b Book) GetLogsByBookStatus(ctx context.Context, r *request.BookGetByStatus) (*response.BookLogsByStatus, error) {
	userID := auth.UserIDFromContext(ctx)
	books, err := b.bookRepo.FindByStatus(ctx, userID, entity.Status(r.Status))
	if err != nil {
		return nil, err
	}
	items := make([]bookWithLogs, 0, len(books))
	for i := range books {
		logs, err := b.logRepo.FindByBookID(ctx, userID, books[i].ID)
		if err != nil {
			return nil, err
		}
		items = append(items, bookWithLogs{Book: books[i], Logs: logs})
	}
	return buildBookLogsByStatusResponse(items), nil
}

func (b Book) GetByID(ctx context.Context, r *request.BookGetByID) (*response.BookGetByID, error) {
	userID := auth.UserIDFromContext(ctx)
	book, err := b.bookRepo.FindByID(ctx, userID, r.BookID)
	if err != nil {
		return nil, err
	}
	return response.NewBookGetByID(book), nil
}

func (b Book) Create(ctx context.Context, r *request.BookCreate) (*response.BookCreate, error) {
	userID := auth.UserIDFromContext(ctx)
	targetCompleteDate := r.TargetCompleteDate.Time()
	if targetCompleteDate.IsZero() {
		targetCompleteDate = time.Now().UTC().AddDate(0, 1, 0)
	}
	book := entity.NewBook(
		userID,
		r.Title,
		r.Author,
		r.Publisher,
		r.ThumbnailUrl,
		r.TotalPages,
		targetCompleteDate,
		r.TargetPagesPerDay,
		entity.Status(r.Status),
		r.ReadPages,
		r.EncounterNote,
	)
	if err := b.bookRepo.Put(ctx, book); err != nil {
		return nil, err
	}
	return response.NewBookCreate(book), nil
}

func (b Book) Update(ctx context.Context, r *request.BookUpdate) (*response.BookUpdate, error) {
	userID := auth.UserIDFromContext(ctx)
	book, err := b.bookRepo.FindByID(ctx, userID, r.BookID)
	if err != nil {
		return nil, err
	}

	// ログがある本は総ページ数を変更できない（バリデーションをパッチ適用より先に行う）
	if r.TotalPages != nil {
		if err := b.ensureNoLogsForBook(ctx, userID, r.BookID); err != nil {
			return nil, err
		}
	}

	var targetCompleteDate *time.Time
	if r.TargetCompleteDate != nil {
		t := r.TargetCompleteDate.Time()
		targetCompleteDate = &t
	}
	var status *entity.Status
	if r.Status != nil {
		s := entity.Status(*r.Status)
		status = &s
	}
	book.ApplyPatch(entity.BookPatch{
		Title:              r.Title,
		Author:             r.Author,
		Publisher:          r.Publisher,
		ThumbnailUrl:       r.ThumbnailUrl,
		TotalPages:         r.TotalPages,
		TargetCompleteDate: targetCompleteDate,
		TargetPagesPerDay:  r.TargetPagesPerDay,
		Status:             status,
		EncounterNote:      r.EncounterNote,
	})
	book.Updated()
	if err := b.bookRepo.Put(ctx, book); err != nil {
		return nil, err
	}
	return response.NewBookUpdate(book), nil
}

func (b Book) Delete(ctx context.Context, r *request.BookDelete) (*response.BookDelete, error) {
	userID := auth.UserIDFromContext(ctx)
	if err := b.logRepo.DeleteByBookID(ctx, userID, r.BookID); err != nil {
		return nil, err
	}
	if err := b.bookRepo.Delete(ctx, userID, r.BookID); err != nil {
		return nil, err
	}
	return response.NewBookDelete(r.BookID), nil
}

func (b Book) ensureNoLogsForBook(ctx context.Context, userID string, bookID int) error {
	logs, err := b.logRepo.FindByBookID(ctx, userID, bookID)
	if err != nil {
		return err
	}
	if len(logs) > 0 {
		return errors.New("totalPages can only be updated when the book has no logs")
	}
	return nil
}

type bookWithLogs struct {
	Book entity.Book
	Logs []entity.Log
}

func buildBookLogsByStatusResponse(items []bookWithLogs) *response.BookLogsByStatus {
	res := &response.BookLogsByStatus{Items: make([]*response.BookWithLogs, 0, len(items))}
	for i := range items {
		cp := items[i].Book
		logs := make([]*entity.Log, 0, len(items[i].Logs))
		for j := range items[i].Logs {
			l := items[i].Logs[j]
			logs = append(logs, &l)
		}
		res.Items = append(res.Items, &response.BookWithLogs{Book: &cp, Logs: logs})
	}
	return res
}
