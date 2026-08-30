package repository

import (
	"context"
	"errors"

	"cloud.google.com/go/datastore"

	"github.com/sora-00/booktracker-api/app/domain/entity"
	domainerr "github.com/sora-00/booktracker-api/app/domain/errors"
	domainrepo "github.com/sora-00/booktracker-api/app/domain/repository"
	dsclient "github.com/sora-00/booktracker-api/app/infra/datastore"
)

const kindLog = "Log"

type logRepo struct{}

func NewLogRepo() domainrepo.LogRepo {
	return &logRepo{}
}

func (r *logRepo) Put(ctx context.Context, log *entity.Log) error {
	if log == nil {
		return errors.New("log is nil")
	}
	ds, err := dsclient.FromContext(ctx)
	if err != nil {
		return err
	}
	var key *datastore.Key
	if log.ID == 0 {
		key = datastore.IncompleteKey(kindLog, nil)
	} else {
		key = datastore.IDKey(kindLog, int64(log.ID), nil)
	}
	key, err = ds.Put(ctx, key, log)
	if err != nil {
		return err
	}
	log.ID = int(key.ID)
	return nil
}

func (r *logRepo) FindAll(ctx context.Context, userID string) ([]entity.Log, error) {
	ds, err := dsclient.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	q := datastore.NewQuery(kindLog).FilterField("userId", "=", userID).Order("readDate")
	var logs []entity.Log
	keys, err := ds.GetAll(ctx, q, &logs)
	if err != nil {
		return nil, err
	}
	for i := range keys {
		logs[i].ID = int(keys[i].ID)
	}
	return logs, nil
}

func (r *logRepo) FindByID(ctx context.Context, userID string, id int) (*entity.Log, error) {
	ds, err := dsclient.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	key := datastore.IDKey(kindLog, int64(id), nil)
	log := &entity.Log{}
	if err := ds.Get(ctx, key, log); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, domainerr.ErrNotFound
		}
		return nil, err
	}
	if log.UserID != userID {
		return nil, domainerr.ErrNotFound
	}
	log.ID = int(key.ID)
	return log, nil
}

func (r *logRepo) FindByBookID(ctx context.Context, userID string, bookID int) ([]entity.Log, error) {
	ds, err := dsclient.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	q := datastore.NewQuery(kindLog).
		FilterField("userId", "=", userID).
		FilterField("bookId", "=", bookID).
		Order("readDate")
	var logs []entity.Log
	keys, err := ds.GetAll(ctx, q, &logs)
	if err != nil {
		return nil, err
	}
	for i := range keys {
		logs[i].ID = int(keys[i].ID)
	}
	return logs, nil
}

func (r *logRepo) DeleteByBookID(ctx context.Context, userID string, bookID int) error {
	logs, err := r.FindByBookID(ctx, userID, bookID)
	if err != nil {
		return err
	}
	ds, err := dsclient.FromContext(ctx)
	if err != nil {
		return err
	}
	for _, log := range logs {
		if err := ds.Delete(ctx, datastore.IDKey(kindLog, int64(log.ID), nil)); err != nil {
			return err
		}
	}
	return nil
}

func (r *logRepo) DeleteByLogID(ctx context.Context, userID string, id int) error {
	ds, err := dsclient.FromContext(ctx)
	if err != nil {
		return err
	}
	if _, err := r.FindByID(ctx, userID, id); err != nil {
		return err
	}
	return ds.Delete(ctx, datastore.IDKey(kindLog, int64(id), nil))
}
