package token

import (
	"context"
	"github.com/api_base/internal/domain/model"
	"github.com/api_base/tool/restclient"
)

const (
	externalApi = "token_api"
)

type Repository struct {
	rc restclient.RestClient
}

func NewRepository(rc restclient.RestClient) *Repository {
	return &Repository{
		rc,
	}
}

func (r *Repository) Get(ctx context.Context, id string) (model.Token, error) {
	result := model.Token{}
	url, err := r.rc.BuildUrl(externalApi, "get_token", id)
	if err != nil {
		return result, err
	}
	err = r.rc.DoGet(ctx, url, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}
