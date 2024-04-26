package user

import (
	"context"

	"github.com/samber/lo"

	"gitlab.seakoi.net/engineer/backend/be-template/constant"
	"gitlab.seakoi.net/engineer/backend/be-template/model/dto/user"
	"gitlab.seakoi.net/engineer/backend/be-template/model/entity"
)

const (
	defaultMemory      = 16 * 1024
	defaultIterations  = 2
	defaultParallelism = 2
	defaultSaltLength  = 16
	defaultKeyLength   = 32
)

func (s *service) Login(ctx context.Context, request *user.LoginRequest) (*user.LoginResponse, error) {
	if err := s.res.Validator().Struct(request); err != nil {
		return nil, err
	}
	user, err := s.userRepository().GetByUsername(ctx, request.Username)
	if err != nil {
		s.res.Logger().Error(ctx, "failed to get user by username", err)
		return nil, constant.ErrSystemMalfunction
	}
	if user == nil {
		return nil, constant.ErrUserNotFound
	}

	compared, err := s.comparePassword(request.Password, user.Password)
	if err != nil {
		s.res.Logger().Error(ctx, "failed to compare password", constant.ErrSystemMalfunction)
		return nil, err
	}
	if !compared {
		return nil, constant.ErrUserPasswordWrong
	}

	return s.generateLoginResponse(ctx, user), nil
}

func (s *service) generateLoginResponse(ctx context.Context, userEntity *entity.User) *user.LoginResponse {
	jwtClaims := userEntity.ToJWTClaims()
	return &user.LoginResponse{
		AccessToken:  lo.Must(s.tokenManger.GenerateAccessToken(ctx, jwtClaims)),
		RefreshToken: lo.Must(s.tokenManger.GenerateRefreshToken(ctx, jwtClaims)),
	}
}

func (s *service) comparePassword(password, encodedHash string) (bool, error) {
	return s.passwordHash.Compare(password, encodedHash)
}
