package httprouter

import (
	"context"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/kucjac/go-rest-sdk/errhandler"
	"github.com/kucjac/go-rest-sdk/forms"
	"github.com/kucjac/go-rest-sdk/handlers"
	"github.com/kucjac/go-rest-sdk/logger"
	"github.com/kucjac/go-rest-sdk/repository"
	"github.com/kucjac/go-rest-sdk/response"
	"net/http"
)

func HttprouterParamGetterFunc(param string, req *http.Request) (paramValue string, err error) {
	params, ok := req.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	if !ok {
		err = errors.New("No httprouter ParamsKey in the request context. Cannot find httprouter.Params")
		return
	}

	paramValue = params.ByName(param)
	return
}

type HttpRouterHandler struct {
	handlers.GenericHandler
}

func New(
	repo repository.Repository,
	errHandler *errhandler.ErrorHandler,
	body response.Responser,
	logs logger.ExtendedLeveledLogger,
) (*HttpRouterHandler, error) {
	generic, err := handlers.New(repo, errHandler, body, logs)
	if err != nil {
		return nil, err
	}
	generic.WithParamGetterFunc(HttprouterParamGetterFunc)

	return &HttpRouterHandler{GenericHandler: *generic}, nil
}

func (h *HttpRouterHandler) New() *HttpRouterHandler {
	handlerCopy := *h
	return &(handlerCopy)
}

func (h *HttpRouterHandler) WithQueryPolicy(policy *forms.BindPolicy) *HttpRouterHandler {
	h.GenericHandler.WithQueryPolicy(policy)
	return h

}

func (h *HttpRouterHandler) WithParamPolicy(policy *forms.ParamPolicy) *HttpRouterHandler {
	h.GenericHandler.WithParamPolicy(policy)
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
