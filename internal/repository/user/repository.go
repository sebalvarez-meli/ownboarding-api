package user

import (
	"context"
	"github.com/api_base/internal/domain/model"
	"github.com/api_base/tool/database"
)

type Repository struct {
	database database.Database
}

func NewRepository(db database.Database) *Repository {
	return &Repository{
		database: db,
	}
}

func (r *Repository) Get(ctx context.Context, id string) (*model.User, error) {
	modelDb := &model.User{
		Id: id,
	}
	return modelDb, nil
}
