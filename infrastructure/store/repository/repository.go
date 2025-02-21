package repository

import (
	"github.com/google/wire"

	"github.com/chaihaobo/be-template/infrastructure/store/repository/user"
)

var ProviderSet = wire.NewSet(
	user.NewRepository,
	New,
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

func New(
	userRepository user.Repository,
) Repository {
	return &repository{
		userRepository: userRepository,
	}
}
