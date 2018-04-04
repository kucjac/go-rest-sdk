package gorillamux

import (
	"github.com/gorilla/mux"
	"github.com/kucjac/go-rest-sdk/errhandler"
	"github.com/kucjac/go-rest-sdk/handlers"
	"github.com/kucjac/go-rest-sdk/logger"
	"github.com/kucjac/go-rest-sdk/repository"
	"github.com/kucjac/go-rest-sdk/response"
	"net/http"
)

func GorillaMuxParamGetterFunc(param string, req *http.Request) (string, error) {
	paramValue := mux.Vars(req)[param]
	return paramValue, nil
}

func New(
	repo repository.Repository,
	errHandler *errhandler.ErrorHandler,
	responseBody response.Responser,
	logs logger.ExtendedLeveledLogger,
) (*handlers.GenericHandler, error) {
	generic, err := handlers.New(repo, errHandler, responseBody, logs)
	if err != nil {
		return nil, err
	}
	generic.WithParamGetterFunc(GorillaMuxParamGetterFunc)
	return generic, nil
}
