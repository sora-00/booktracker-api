package repository

import (
	"context"
	"errors"

	"cloud.google.com/go/datastore"

	"github.com/sora-00/booktracker-api/app/domain/entity"
	dsclient "github.com/sora-00/booktracker-api/app/infra/datastore"
)

// BookRepo は本の永続化のインターフェース。
type BookRepo interface {
	Create(ctx context.Context, book *entity.Book) error
	FindAll(ctx context.Context) ([]entity.Book, error)
	FindByID(ctx context.Context, id int) (*entity.Book, error)
	Delete(ctx context.Context, id int) error
}

// ErrNotFound は対象が存在しないときに返す。controller で 404 に変換する。
var ErrNotFound = errors.New("not found")

const kindBook = "Book"

type bookRepo struct{}

func NewBookRepo() BookRepo {
	return &bookRepo{}
}

var errNoDatastore = errors.New("datastore client not found in context")

func (r *bookRepo) ds(ctx context.Context) (*datastore.Client, error) {
	ds, ok := dsclient.FromContext(ctx)
	if !ok {
		return nil, errNoDatastore
	}
	return ds, nil
}

func (r *bookRepo) Create(ctx context.Context, book *entity.Book) error {
	ds, err := r.ds(ctx)
	if err != nil {
		return err
	}
	key := datastore.IncompleteKey(kindBook, nil)
	key, err = ds.Put(ctx, key, book)
	if err != nil {
		return err
	}
	book.ID = int(key.ID)
	return nil
}

func (r *bookRepo) FindAll(ctx context.Context) ([]entity.Book, error) {
	ds, err := r.ds(ctx)
	if err != nil {
		return nil, err
	}
	q := datastore.NewQuery(kindBook).Order("createdAt")
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

func (r *bookRepo) FindByID(ctx context.Context, id int) (*entity.Book, error) {
	ds, err := r.ds(ctx)
	if err != nil {
		return nil, err
	}
	key := datastore.IDKey(kindBook, int64(id), nil)
	book := &entity.Book{}
	if err := ds.Get(ctx, key, book); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, ErrNotFound
		}
		return nil, err
	}
	book.ID = int(key.ID)
	return book, nil
}

func (r *bookRepo) Delete(ctx context.Context, id int) error {
	ds, err := r.ds(ctx)
	if err != nil {
		return err
	}
	if _, err := r.FindByID(ctx, id); err != nil {
		return err
	}
	key := datastore.IDKey(kindBook, int64(id), nil)
	return ds.Delete(ctx, key)
}
