package service

import (
	"context"
	"errors"
	"time"

	"github.com/sora-00/booktracker-api/app/domain/entity"
	"github.com/sora-00/booktracker-api/app/domain/repository"
)

type LogSvc struct {
	logRepo repository.LogRepo
}

func NewLogService(logRepo repository.LogRepo) *LogSvc {
	return &LogSvc{logRepo: logRepo}
}

func (s *LogSvc) CreateLog(ctx context.Context, log *entity.Log) (*entity.Log, error) {
	if log == nil {
		return nil, errors.New("log is required")
	}
	log.ComputePagesRead()
	if err := s.logRepo.Put(ctx, log); err != nil {
		return nil, err
	}
	return log, nil
}

func (s *LogSvc) UpdateLog(ctx context.Context, log *entity.Log) (*entity.Log, error) {
	if log == nil {
		return nil, errors.New("log is required")
	}
	log.ComputePagesRead()
	if err := s.logRepo.Put(ctx, log); err != nil {
		return nil, err
	}
	return log, nil
}

// ValidateCreate は新規ログ登録のドメインルールを検証する。
func (s *LogSvc) ValidateCreate(readDate time.Time, startPage, endPage, pagesAll int, existing []entity.Log) error {
	if !readDate.After(latestReadDateOf(existing)) {
		return errors.New("readDate must be after the latest log's readDate for this book")
	}
	if startPage < maxEndPageOf(existing) {
		return errors.New("startPage must be >= previous log's endPage")
	}
	return validateEndPage(endPage, pagesAll)
}

// ValidateUpdate はログ更新のドメインルールを検証する。
func (s *LogSvc) ValidateUpdate(updated *entity.Log, pagesAll int, existing []entity.Log) error {
	if err := validateEndPage(updated.EndPage, pagesAll); err != nil {
		return err
	}
	prevEnd := maxEndPageBeforeDate(existing, updated.ID, updated.ReadDate)
	if updated.StartPage < prevEnd {
		return errors.New("startPage must be >= previous log's endPage")
	}
	if hasFutureOverlap(existing, updated.ID, updated.ReadDate, updated.EndPage) {
		return errors.New("log order constraint violated")
	}
	return nil
}

// MaxEndPageOf は readPages 更新に使うログ群の最大 endPage を返す。
func MaxEndPageOf(logs []entity.Log) int {
	var max int
	for i := range logs {
		if logs[i].EndPage > max {
			max = logs[i].EndPage
		}
	}
	return max
}

func validateEndPage(endPage, pagesAll int) error {
	if endPage > pagesAll {
		return errors.New("endPage must not exceed book's pagesAll")
	}
	return nil
}

func latestReadDateOf(logs []entity.Log) time.Time {
	var latest time.Time
	for i := range logs {
		if logs[i].ReadDate.After(latest) {
			latest = logs[i].ReadDate
		}
	}
	return latest
}

func maxEndPageOf(logs []entity.Log) int {
	var max int
	for i := range logs {
		if logs[i].EndPage > max {
			max = logs[i].EndPage
		}
	}
	return max
}

func maxEndPageBeforeDate(logs []entity.Log, selfID int, d time.Time) int {
	var prevEnd int
	for i := range logs {
		if logs[i].ID == selfID {
			continue
		}
		if logs[i].ReadDate.Before(d) && logs[i].EndPage > prevEnd {
			prevEnd = logs[i].EndPage
		}
	}
	return prevEnd
}

func hasFutureOverlap(logs []entity.Log, selfID int, d time.Time, endPage int) bool {
	for i := range logs {
		if logs[i].ID == selfID {
			continue
		}
		if logs[i].ReadDate.After(d) && logs[i].StartPage < endPage {
			return true
		}
	}
	return false
}
