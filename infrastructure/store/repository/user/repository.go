package user

import (
	"context"

	"gitlab.seakoi.net/engineer/backend/be-template/infrastructure/store/client"
	"gitlab.seakoi.net/engineer/backend/be-template/model/entity"
)

type (
	Repository interface {
		GetByUsername(ctx context.Context, username string) (*entity.User, error)
	}
	repository struct {
		client client.Client
	}
)

func (r repository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var data entity.User
	if result := r.client.DB(ctx).Where("username = ?", username).Find(&data); result.Error != nil || result.RowsAffected == 0 {
		return nil, result.Error
	}
	return &data, nil
}

func NewRepository(client client.Client) Repository {
	return &repository{
		client: client,
	}
}
