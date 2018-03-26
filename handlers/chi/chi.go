package chihandler

import (
	"github.com/go-chi/chi"
	"github.com/kucjac/go-rest-sdk/errhandler"
	"github.com/kucjac/go-rest-sdk/forms"
	"github.com/kucjac/go-rest-sdk/handlers"
	"github.com/kucjac/go-rest-sdk/refutils"
	"github.com/kucjac/go-rest-sdk/repository"
	"github.com/kucjac/go-rest-sdk/response"
	"net/http"
)

func ChiSetIDFunc(
	req *http.Request,
	model interface{},
) (err error) {
	modelName := refutils.ModelName(model)
	modelID := chi.URLParam(req, modelName)
	if modelID == "" {
		err = handlers.ErrIncorrectModel
		return err
	}
	err = forms.SetID(model, modelID)
	if err != nil {
		return err
	}
	return nil
}

func New(
	repo repository.Repository,
	errHandler *errhandler.ErrorHandler,
	body response.Responser,
) (*handlers.GenericHandler, error) {
	h, err := handlers.New(repo, errHandler, body)
	if err != nil {
		return nil, err
	}
	return h.WithSetIDFunc(ChiSetIDFunc), nil
}
