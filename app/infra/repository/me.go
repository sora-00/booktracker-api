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

const kindMe = "Me"

type meRepo struct{}

func NewMeRepo() domainrepo.MeRepo {
	return &meRepo{}
}

func (r *meRepo) Put(ctx context.Context, m *entity.Me) error {
	if m == nil {
		return errors.New("me is nil")
	}
	if m.ID == "" {
		return errors.New("me id is required")
	}
	ds, err := dsclient.FromContext(ctx)
	if err != nil {
		return err
	}
	_, err = ds.Put(ctx, datastore.NameKey(kindMe, m.ID, nil), m)
	return err
}

func (r *meRepo) Get(ctx context.Context, uid string) (*entity.Me, error) {
	ds, err := dsclient.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	m := &entity.Me{}
	if err := ds.Get(ctx, datastore.NameKey(kindMe, uid, nil), m); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, domainerr.ErrNotFound
		}
		return nil, err
	}
	m.ID = uid
	return m, nil
}

func (r *meRepo) Delete(ctx context.Context, uid string) error {
	ds, err := dsclient.FromContext(ctx)
	if err != nil {
		return err
	}
	if err := ds.Delete(ctx, datastore.NameKey(kindMe, uid, nil)); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return domainerr.ErrNotFound
		}
		return err
	}
	return nil
}
