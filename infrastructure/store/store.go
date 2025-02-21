package store

import (
	"github.com/google/wire"

	"github.com/chaihaobo/be-template/infrastructure/store/client"
	"github.com/chaihaobo/be-template/infrastructure/store/repository"
)

var ProviderSet = wire.NewSet(
	client.ProviderSet,
	repository.ProviderSet,
	New,
)

type (
	Store interface {
		Client() client.Client
		Repository() repository.Repository
	}
	store struct {
		client     client.Client
		repository repository.Repository
	}
)

func (s *store) Repository() repository.Repository {
	return s.repository
}

func (s *store) Client() client.Client {
	return s.client
}

func New(client client.Client, repository repository.Repository) (Store, error) {
	return &store{
		client:     client,
		repository: repository,
	}, nil
}
