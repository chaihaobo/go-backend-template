package infrastructure

import (
	"github.com/google/wire"

	"github.com/chaihaobo/be-template/infrastructure/cache"
	"github.com/chaihaobo/be-template/infrastructure/store"
)

var ProviderSet = wire.NewSet(
	cache.ProviderSet,
	store.ProviderSet,
	New,
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

func New(store store.Store, cacheClient cache.Client) Infrastructure {
	return &infrastructure{
		store: store,
		cache: cacheClient,
	}
}
