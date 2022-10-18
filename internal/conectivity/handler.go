package conectivity

import (
	"context"
	"errors"
	"github.com/api_base/internal/conectivity/response"
	"github.com/api_base/internal/domain/model"
	"github.com/go-chi/chi"
	"net/http"
)

type HandlerFunc interface {
	Get(w http.ResponseWriter, r *http.Request)
}

type Service interface {
	Get(ctx context.Context, id string) (*model.User, error)
}

type handler struct {
	service Service
}

func NewHandlerFunc(srv Service) HandlerFunc {
	return &handler{service: srv}
}

func (h handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Write(w, errors.New("invalid_id"), http.StatusBadRequest)
	}
	user, err := h.service.Get(ctx, id)
	if err != nil {
		response.Write(w, err, http.StatusInternalServerError)
	}
	response.Write(w, user, http.StatusOK)
}
