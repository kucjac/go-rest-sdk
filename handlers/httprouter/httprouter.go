package httprouter

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/kucjac/go-rest-sdk/errhandler"
	"github.com/kucjac/go-rest-sdk/forms"
	"github.com/kucjac/go-rest-sdk/handlers"
	"github.com/kucjac/go-rest-sdk/refutils"
	"github.com/kucjac/go-rest-sdk/repository"
	"github.com/kucjac/go-rest-sdk/response"
	"net/http"
)

func HttpRouterSetIDFunc(req *http.Request, model interface{}) error {

	params, ok := req.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	if !ok {
		return handlers.ErrIncorrectCustomContext
	}
	name := refutils.ModelName(model)
	paramID := params.ByName(name)
	if err := forms.SetID(model, paramID); err != nil {
		return err
	}
	return nil
}

type HttpRouterHandler struct {
	handlers.GenericHandler
}

func New(
	repo repository.Repository,
	errHandler *errhandler.ErrorHandler,
	body response.Responser,
) (*HttpRouterHandler, error) {
	h, err := handlers.New(repo, errHandler, body)
	if err != nil {
		return nil, err
	}
	h.WithSetIDFunc(HttpRouterSetIDFunc)
	return &HttpRouterHandler{GenericHandler: *h}, nil
}

func (h *HttpRouterHandler) New() *HttpRouterHandler {
	return &(*h)
}

func (h *HttpRouterHandler) WithQueryPolicy(policy *forms.Policy) *HttpRouterHandler {
	h.GenericHandler.WithQueryPolicy(policy)
	return h
}

func (h *HttpRouterHandler) WithJSONPolicy(policy *forms.Policy) *HttpRouterHandler {
	h.GenericHandler.WithJSONPolicy(policy)
	return h
}

func (h *HttpRouterHandler) WithListPolicy(policy *forms.ListPolicy) *HttpRouterHandler {
	h.GenericHandler.WithListPolicy(policy)
	return h
}

func (h *HttpRouterHandler) Create(model interface{}) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
		contextWithParams(req, p)
		h.GenericHandler.Create(model).ServeHTTP(rw, req)
	}
}

func (h *HttpRouterHandler) Get(model interface{}) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
		contextWithParams(req, p)
		h.GenericHandler.Get(model).ServeHTTP(rw, req)
	}
}

func (h *HttpRouterHandler) List(model interface{}) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
		contextWithParams(req, p)
		h.GenericHandler.List(model).ServeHTTP(rw, req)
	}
}

func (h *HttpRouterHandler) Update(model interface{}) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
		contextWithParams(req, p)
		h.GenericHandler.Update(model).ServeHTTP(rw, req)
	}
}

func (h *HttpRouterHandler) Patch(model interface{}) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
		contextWithParams(req, p)
		h.GenericHandler.Patch(model).ServeHTTP(rw, req)
	}
}

func (h *HttpRouterHandler) Delete(model interface{}) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
		contextWithParams(req, p)
		h.GenericHandler.Delete(model).ServeHTTP(rw, req)
	}
}

func contextWithParams(req *http.Request, p httprouter.Params) {
	ctx := req.Context()
	ctx = context.WithValue(ctx, httprouter.ParamsKey, p)
	req.WithContext(ctx)
}
