package chihandler

import (
	"github.com/go-chi/chi"
	"github.com/kucjac/go-rest-sdk/errhandler"
	"github.com/kucjac/go-rest-sdk/handlers"
	"github.com/kucjac/go-rest-sdk/logger"
	"github.com/kucjac/go-rest-sdk/repository"
	"github.com/kucjac/go-rest-sdk/response"
	"net/http"
)

func ChiParamGetterFunc(param string, req *http.Request) (string, error) {
	paramValue := chi.URLParam(req, param)
	return paramValue, nil
}

func New(
	repo repository.Repository,
	errHandler *errhandler.ErrorHandler,
	body response.Responser,
	logs logger.ExtendedLeveledLogger,
) (*handlers.GenericHandler, error) {
	generic, err := handlers.New(repo, errHandler, body, logs)
	if err != nil {
		return nil, err
	}
	generic.WithParamGetterFunc(ChiParamGetterFunc)
	return generic, nil
}
