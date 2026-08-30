package repository

import "context"

type BookThumbnailRepo interface {
	Save(ctx context.Context, name string, data []byte) error
	Load(ctx context.Context, name string) ([]byte, error)
}
