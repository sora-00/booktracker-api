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

const kindBook = "Book"

type bookRepo struct{}

func NewBookRepo() domainrepo.BookRepo {
	return &bookRepo{}
}

func (r *bookRepo) Put(ctx context.Context, book *entity.Book) error {
	if book == nil {
		return errors.New("book is nil")
	}
	ds, err := dsclient.FromContext(ctx)
	if err != nil {
		return err
	}
	var key *datastore.Key
	if book.ID == 0 {
		key = datastore.IncompleteKey(kindBook, nil)
	} else {
		key = datastore.IDKey(kindBook, int64(book.ID), nil)
	}
	key, err = ds.Put(ctx, key, book)
	if err != nil {
		return err
	}
	book.ID = int(key.ID)
	return nil
}

func (r *bookRepo) FindAll(ctx context.Context, userID string) ([]entity.Book, error) {
	ds, err := dsclient.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	q := datastore.NewQuery(kindBook).FilterField("userId", "=", userID).Order("-createdAt")
	var books []entity.Book
	keys, err := ds.GetAll(ctx, q, &books)
	if err != nil {
		return nil, err
	}
	for i := range keys {
		books[i].ID = int(keys[i].ID)
	}
	return books, nil
}

func (r *bookRepo) FindByStatus(ctx context.Context, userID string, status entity.Status) ([]entity.Book, error) {
	ds, err := dsclient.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	q := datastore.NewQuery(kindBook).
		FilterField("userId", "=", userID).
		FilterField("status", "=", string(status)).
		Order("-createdAt")
	var books []entity.Book
	keys, err := ds.GetAll(ctx, q, &books)
	if err != nil {
		return nil, err
	}
	for i := range keys {
		books[i].ID = int(keys[i].ID)
	}
	return books, nil
}

func (r *bookRepo) FindByID(ctx context.Context, userID string, id int) (*entity.Book, error) {
	ds, err := dsclient.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	key := datastore.IDKey(kindBook, int64(id), nil)
	book := &entity.Book{}
	if err := ds.Get(ctx, key, book); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, domainerr.ErrNotFound
		}
		return nil, err
	}
	if book.UserID != userID {
		return nil, domainerr.ErrNotFound
	}
	book.ID = int(key.ID)
	return book, nil
}

func (r *bookRepo) Delete(ctx context.Context, userID string, id int) error {
	ds, err := dsclient.FromContext(ctx)
	if err != nil {
		return err
	}
	if _, err := r.FindByID(ctx, userID, id); err != nil {
		return err
	}
	return ds.Delete(ctx, datastore.IDKey(kindBook, int64(id), nil))
}
