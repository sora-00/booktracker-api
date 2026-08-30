package usecase

import (
	"context"

	"github.com/sora-00/booktracker-api/app/domain/entity"
	"github.com/sora-00/booktracker-api/app/domain/repository"
	"github.com/sora-00/booktracker-api/app/infra/auth"
)

type Me struct {
	meRepo repository.MeRepo
}

func NewMe(meRepo repository.MeRepo) *Me {
	return &Me{meRepo: meRepo}
}

func (u Me) Get(ctx context.Context) (*entity.Me, error) {
	return u.meRepo.Get(ctx, auth.UserIDFromContext(ctx))
}

func (u Me) Delete(ctx context.Context) error {
	return u.meRepo.Delete(ctx, auth.UserIDFromContext(ctx))
}
