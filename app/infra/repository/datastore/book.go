package datastore

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/datastore"

	"github.com/sora-00/booktracker-api/app/domain/entity"
	"github.com/sora-00/booktracker-api/app/domain/repository"
	dsclient "github.com/sora-00/booktracker-api/pkg/datastore"
)

const kindBook = "Book"

// bookEntity は Datastore に保存するための構造体（ID は Key で持つ）
type bookEntity struct {
	Title              string    `datastore:"title"`
	Author             string    `datastore:"author"`
	TotalPages         int       `datastore:"totalPages"`
	Publisher          string    `datastore:"publisher"`
	ThumbnailUrl       string    `datastore:"thumbnailUrl"`
	Status             string    `datastore:"status"`
	TargetCompleteDate time.Time `datastore:"targetCompleteDate"`
	CreatedAt          time.Time `datastore:"createdAt"`
	UpdatedAt          time.Time `datastore:"updatedAt"`
}

type bookRepo struct{}

// NewBookRepo は BookRepo の Cloud Datastore 実装を返す。
// 各メソッドでは context から ds を取得する（middleware で WithContext されている前提）。
func NewBookRepo() repository.BookRepo {
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

func toEntity(key *datastore.Key, e *bookEntity) *entity.Book {
	if key == nil {
		return nil
	}
	return &entity.Book{
		ID:                 int(key.ID),
		Title:              e.Title,
		Author:             e.Author,
		TotalPages:         e.TotalPages,
		Publisher:          e.Publisher,
		ThumbnailUrl:       e.ThumbnailUrl,
		Status:             entity.Status(e.Status),
		TargetCompleteDate: e.TargetCompleteDate,
		CreatedAt:          e.CreatedAt,
		UpdatedAt:          e.UpdatedAt,
	}
}

func toBookEntity(b *entity.Book) *bookEntity {
	return &bookEntity{
		Title:              b.Title,
		Author:             b.Author,
		TotalPages:         b.TotalPages,
		Publisher:          b.Publisher,
		ThumbnailUrl:       b.ThumbnailUrl,
		Status:             string(b.Status),
		TargetCompleteDate: b.TargetCompleteDate,
		CreatedAt:          b.CreatedAt,
		UpdatedAt:          b.UpdatedAt,
	}
}

func (r *bookRepo) Create(ctx context.Context, book *entity.Book) error {
	ds, err := r.ds(ctx)
	if err != nil {
		return err
	}
	key := datastore.IncompleteKey(kindBook, nil)
	e := toBookEntity(book)
	key, err = ds.Put(ctx, key, e)
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
	var entities []bookEntity
	keys, err := ds.GetAll(ctx, q, &entities)
	if err != nil {
		return nil, err
	}
	books := make([]entity.Book, 0, len(keys))
	for i := range keys {
		books = append(books, *toEntity(keys[i], &entities[i]))
	}
	return books, nil
}

func (r *bookRepo) FindByID(ctx context.Context, id int) (*entity.Book, error) {
	ds, err := r.ds(ctx)
	if err != nil {
		return nil, err
	}
	key := datastore.IDKey(kindBook, int64(id), nil)
	var e bookEntity
	if err := ds.Get(ctx, key, &e); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return toEntity(key, &e), nil
}

func (r *bookRepo) Delete(ctx context.Context, id int) error {
	ds, err := r.ds(ctx)
	if err != nil {
		return err
	}
	if _, err := r.FindByID(ctx, id); err != nil {
		return err // 存在しなければ ErrNotFound
	}
	key := datastore.IDKey(kindBook, int64(id), nil)
	return ds.Delete(ctx, key)
}
