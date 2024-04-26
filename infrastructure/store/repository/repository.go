package repository

import (
	"gitlab.seakoi.net/engineer/backend/be-template/infrastructure/store/client"
	"gitlab.seakoi.net/engineer/backend/be-template/infrastructure/store/repository/user"
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
