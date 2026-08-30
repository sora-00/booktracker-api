package usecase

import (
	"context"
	"errors"

	"github.com/sora-00/booktracker-api/app/domain/entity"
	"github.com/sora-00/booktracker-api/app/domain/repository"
	"github.com/sora-00/booktracker-api/app/domain/service"
	"github.com/sora-00/booktracker-api/app/infra/auth"
	"github.com/sora-00/booktracker-api/app/usecase/request"
	"github.com/sora-00/booktracker-api/app/usecase/response"
)

type Log struct {
	logRepo    repository.LogRepo
	bookRepo   repository.BookRepo
	logService *service.LogSvc
}

func NewLog(logRepo repository.LogRepo, bookRepo repository.BookRepo, logService *service.LogSvc) *Log {
	return &Log{logRepo: logRepo, bookRepo: bookRepo, logService: logService}
}

func (u Log) Get(ctx context.Context, _ *request.LogGet) (*response.LogGet, error) {
	userID := auth.UserIDFromContext(ctx)
	logs, err := u.logRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	return response.NewLogGet(logs), nil
}

func (u Log) GetByID(ctx context.Context, r *request.LogGetByID) (*response.LogGetByID, error) {
	userID := auth.UserIDFromContext(ctx)
	log, err := u.logRepo.FindByID(ctx, userID, r.LogID)
	if err != nil {
		return nil, err
	}
	return response.NewLogGetByID(log), nil
}

func (u Log) GetByBookID(ctx context.Context, r *request.LogGetByBookID) (*response.LogGetByBookID, error) {
	userID := auth.UserIDFromContext(ctx)
	logs, err := u.logRepo.FindByBookID(ctx, userID, r.BookID)
	if err != nil {
		return nil, err
	}
	return response.NewLogGetByBookID(logs), nil
}

func (u Log) Create(ctx context.Context, r *request.LogCreate) (*response.LogCreate, error) {
	userID := auth.UserIDFromContext(ctx)

	book, err := u.bookRepo.FindByID(ctx, userID, r.BookID)
	if err != nil {
		return nil, err
	}
	readDate, err := request.ParseReadDate(r.ReadDate)
	if err != nil {
		return nil, err
	}
	existing, err := u.logRepo.FindByBookID(ctx, userID, r.BookID)
	if err != nil {
		return nil, err
	}
	if err := u.logService.ValidateCreate(readDate, r.StartPage, r.EndPage, book.TotalPages, existing); err != nil {
		return nil, err
	}

	newLog := entity.NewLog(userID, r.BookID, readDate, r.StartPage, r.EndPage, r.Memo)
	created, err := u.logService.CreateLog(ctx, newLog)
	if err != nil {
		return nil, err
	}
	if err := u.syncBookReadPages(ctx, userID, r.BookID); err != nil {
		return nil, err
	}
	return response.NewLogCreate(created), nil
}

func (u Log) Update(ctx context.Context, r *request.LogUpdate) (*response.LogUpdate, error) {
	userID := auth.UserIDFromContext(ctx)

	target, err := u.logRepo.FindByID(ctx, userID, r.LogID)
	if err != nil {
		return nil, err
	}
	book, err := u.bookRepo.FindByID(ctx, userID, target.BookID)
	if err != nil {
		return nil, err
	}

	draft, err := applyLogPatch(target, r)
	if err != nil {
		return nil, err
	}
	existing, err := u.logRepo.FindByBookID(ctx, userID, draft.BookID)
	if err != nil {
		return nil, err
	}
	if err := u.logService.ValidateUpdate(draft, book.TotalPages, existing); err != nil {
		return nil, err
	}

	draft.Updated()
	updated, err := u.logService.UpdateLog(ctx, draft)
	if err != nil {
		return nil, err
	}
	if err := u.syncBookReadPages(ctx, userID, draft.BookID); err != nil {
		return nil, err
	}
	return response.NewLogUpdate(updated), nil
}

func (u Log) Delete(ctx context.Context, r *request.LogDelete) error {
	userID := auth.UserIDFromContext(ctx)
	log, err := u.logRepo.FindByID(ctx, userID, r.LogID)
	if err != nil {
		return err
	}
	bookID := log.BookID
	if err := u.logRepo.DeleteByLogID(ctx, userID, r.LogID); err != nil {
		return err
	}
	return u.syncBookReadPages(ctx, userID, bookID)
}

func (u Log) DeleteByBookID(ctx context.Context, r *request.LogDeleteByBookID) error {
	userID := auth.UserIDFromContext(ctx)
	if err := u.logRepo.DeleteByBookID(ctx, userID, r.BookID); err != nil {
		return err
	}
	return u.syncBookReadPages(ctx, userID, r.BookID)
}

// syncBookReadPages はログの最大 endPage を book.ReadPages に反映する。
func (u Log) syncBookReadPages(ctx context.Context, userID string, bookID int) error {
	logs, err := u.logRepo.FindByBookID(ctx, userID, bookID)
	if err != nil {
		return err
	}
	book, err := u.bookRepo.FindByID(ctx, userID, bookID)
	if err != nil {
		return err
	}
	book.ReadPages = service.MaxEndPageOf(logs)
	book.Updated()
	return u.bookRepo.Put(ctx, book)
}

// applyLogPatch は既存ログにリクエストの差分を適用した新しい Log を返す。
func applyLogPatch(target *entity.Log, r *request.LogUpdate) (*entity.Log, error) {
	draft := *target
	var patch entity.LogPatch

	if r.ReadDate != nil {
		rd, err := request.ParseReadDate(*r.ReadDate)
		if err != nil {
			return nil, err
		}
		patch.ReadDate = &rd
	}
	patch.StartPage = r.StartPage
	patch.EndPage = r.EndPage
	patch.Memo = r.Memo
	draft.ApplyPatch(patch)

	if draft.StartPage >= draft.EndPage {
		return nil, errors.New("startPage must be less than endPage")
	}
	return &draft, nil
}
