package domain

import (
	"context"
	"github.com/api_base/config"
	"github.com/api_base/internal/domain/model"
	"github.com/api_base/internal/repository/token"
	"github.com/api_base/internal/repository/user"
	"github.com/api_base/tool/database"
	"github.com/api_base/tool/restclient"
	"log"
)

type Container struct {
	UserRepo  UserRepository
	TokenRepo TokenRepository
}

type UserRepository interface {
	Get(ctx context.Context, id string) (*model.User, error)
}

type TokenRepository interface {
	Get(ctx context.Context, id string) (model.Token, error)
}

func NewContainer(config config.Config) Container {
	db, err := database.NewRepository(config.Database)
	if err != nil {
		log.Fatal("initialize database fail: ", err)
	}
	rc, err := restclient.NewRestClient(config.RestClient)
	if err != nil {
		log.Fatal("initialize rest_client fail: ", err)
	}
	return Container{
		UserRepo:  user.NewRepository(db),
		TokenRepo: token.NewRepository(rc),
	}
}
