package infrastructure

import (
	"github.com/chaihaobo/be-template/infrastructure/cache"
	"github.com/chaihaobo/be-template/infrastructure/store"
	"github.com/chaihaobo/be-template/resource"
)

type (
	Infrastructure interface {
		Store() store.Store
		Cache() cache.Client
		Close() error
	}

	infrastructure struct {
		store store.Store
		cache cache.Client
	}
)

func (i *infrastructure) Close() error {
	closeFuncs := []func() error{
		i.store.Client().Close,
		i.cache.Close,
	}
	for _, closeFun := range closeFuncs {
		if err := closeFun(); err != nil {
			return err
		}
	}
	return nil
}

func (i *infrastructure) Cache() cache.Client {
	return i.cache
}

func (i *infrastructure) Store() store.Store {
	return i.store
}

func New(res resource.Resource) (Infrastructure, error) {
	store, err := store.New(res)
	if err != nil {
		return nil, err
	}

	cacheClient, err := cache.NewClient(res)
	if err != nil {
		return nil, err
	}

	return &infrastructure{
		store: store,
		cache: cacheClient,
	}, nil
}
