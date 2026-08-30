package repository

import (
	"context"

	"github.com/sora-00/booktracker-api/app/domain/entity"
)

type MeRepo interface {
	Put(ctx context.Context, m *entity.Me) error
	Get(ctx context.Context, uid string) (*entity.Me, error)
	Delete(ctx context.Context, uid string) error
}
