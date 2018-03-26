package gorillamux

import (
	"github.com/gorilla/mux"
	"github.com/kucjac/go-rest-sdk/errhandler"
	"github.com/kucjac/go-rest-sdk/forms"
	"github.com/kucjac/go-rest-sdk/handlers"
	"github.com/kucjac/go-rest-sdk/refutils"
	"github.com/kucjac/go-rest-sdk/repository"
	"github.com/kucjac/go-rest-sdk/response"
	"net/http"
)

func GorillaMuxIDSetFunc(
	req *http.Request,
	model interface{},
) (err error) {
	vars := mux.Vars(req)

	modelName := refutils.ModelName(model)
	modelID, ok := vars[modelName]
	if !ok {
		err = handlers.ErrIncorrectModel
		return err
	}

	if err = forms.SetID(model, modelID); err != nil {
		return err
	}
	return nil
}

func New(
	repo repository.Repository,
	errHandler *errhandler.ErrorHandler,
	responseBody response.Responser,
) (*handlers.GenericHandler, error) {
	h, err := handlers.New(repo, errHandler, responseBody)
	if err != nil {
		return nil, err
	}
	return h.WithSetIDFunc(GorillaMuxIDSetFunc), nil
}
