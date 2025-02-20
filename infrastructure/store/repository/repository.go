package repository

import (
	"github.com/chaihaobo/be-template/infrastructure/store/client"
	"github.com/chaihaobo/be-template/infrastructure/store/repository/user"
)

type (
	Repository interface {
		User() user.Repository
	}
	repository struct {
		userRepository user.Repository
	}
)

func (r *repository) User() user.Repository {
	return r.userRepository
}

func New(client client.Client) Repository {
	return &repository{
		userRepository: user.NewRepository(client),
	}
}
