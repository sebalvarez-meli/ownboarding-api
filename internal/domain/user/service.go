package user

import (
	"context"
	"github.com/api_base/internal/domain"
	"github.com/api_base/internal/domain/model"
)

type Service interface {
	Get(ctx context.Context, id int64) (*model.User, error)
}

type service struct {
	container domain.Container
}

func NewService(container domain.Container) Service {
	return &service{
		container,
	}
}

func (s service) Get(ctx context.Context, id int64) (*model.User, error) {
	token, err := s.container.TokenRepo.Get(ctx, string(id))
	if err != nil {
		return nil, err
	}
	user, err := s.container.UserRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	user.Token = token
	return user, nil
}
