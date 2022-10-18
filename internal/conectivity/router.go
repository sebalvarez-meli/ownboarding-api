package conectivity

import "github.com/go-chi/chi"

type RouterHandler interface {
	Handler() *chi.Mux
}

type routerHandler struct {
	handlerFunc HandlerFunc
}

func NewRouterHandler(hdlFunc HandlerFunc) RouterHandler {
	return &routerHandler{
		handlerFunc: hdlFunc,
	}
}

func (rh routerHandler) Handler() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/get/{id}", rh.handlerFunc.Get)
	return r
}
