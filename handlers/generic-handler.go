package handlers

import (
	"encoding/json"
	"errors"
	"github.com/kucjac/go-rest-sdk/dberrors"
	"github.com/kucjac/go-rest-sdk/errhandler"
	"github.com/kucjac/go-rest-sdk/forms"
	"github.com/kucjac/go-rest-sdk/refutils"
	"github.com/kucjac/go-rest-sdk/repository"
	"github.com/kucjac/go-rest-sdk/response"
	"github.com/kucjac/go-rest-sdk/resterrors"
	"net/http"
)

var (
	ErrIncorrectModel         = errors.New("Incorrect model route path provided.")
	ErrIncorrectCustomContext = errors.New("Incorrect custom context type.")
)

type GenericHandler struct {
	repo          repository.Repository
	errHandler    *errhandler.ErrorHandler
	queryPolicy   *forms.Policy
	jsonPolicy    *forms.Policy
	listPolicy    *forms.ListPolicy
	responseBody  response.Responser
	idSetFunc     SetIDFunc
	customContext interface{}
}

type SetIDFunc func(req *http.Request, model interface{}) error

// New creates JSONHandler for given
func New(repo repository.Repository,
	errHandler *errhandler.ErrorHandler,
	responseBody response.Responser,
) (*GenericHandler, error) {
	if repo == nil || errHandler == nil {
		return nil, errors.New("repository and errorHandler cannot be nil.")
	}
	if responseBody == nil {
		responseBody = &response.DefaultBody{}
	}
	chiHandler := &GenericHandler{
		repo:         repo,
		errHandler:   errHandler,
		responseBody: responseBody,
	}
	return chiHandler, nil
}

func (c *GenericHandler) New() *GenericHandler {
	h := *c
	return &h
}

func (c *GenericHandler) WithQueryPolicy(policy *forms.Policy) *GenericHandler {
	c.queryPolicy = policy
	return c
}

func (c *GenericHandler) WithJSONPolicy(policy *forms.Policy) *GenericHandler {
	c.jsonPolicy = policy
	return c
}

func (c *GenericHandler) WithListPolicy(policy *forms.ListPolicy) *GenericHandler {
	c.listPolicy = policy
	return c
}

func (c *GenericHandler) WithSetIDFunc(
	customIDFunc SetIDFunc,
) *GenericHandler {
	c.idSetFunc = customIDFunc
	return c
}

// Create is a chiHandler HandlerFunc for creating new restful model records
func (c *GenericHandler) Create(model interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		var status int
		obj := refutils.ObjOfPtrType(model)
		err := forms.BindJSON(req, obj, c.jsonPolicy)
		if err != nil {
			restErr := resterrors.ErrInvalidJSONDocument.New()
			restErr.AddDetailInfo(err.Error())
			status = 400
			c.JSON(rw, req, status, c.getResponseBodyErr(status, restErr))
			return
		}

		dbErr := c.repo.Create(obj)
		if dbErr != nil {
			c.handleDBError(rw, req, dbErr)
			return
		}
		status = http.StatusCreated
		c.JSON(rw, req, status, c.getResponseBodyContent(status, obj))
		return
	}
}

func (c *GenericHandler) Get(model interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {

		obj := refutils.ObjOfPtrType(model)

		err := c.idSetFunc(req, model)
		if err != nil {
			restErr := resterrors.ErrInternalError.New()
			c.JSON(rw, req, 500, c.getResponseBodyErr(500, restErr))
			return
		}

		result, dbErr := c.repo.Get(obj)
		if dbErr != nil {
			c.handleDBError(rw, req, dbErr)
			return
		}

		c.JSON(rw, req, 200, c.getResponseBodyContent(200, result))
	}
}

func (c *GenericHandler) List(model interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		// New 'model' entity
		obj := refutils.ObjOfPtrType(model)

		err := forms.BindQuery(req, obj, c.queryPolicy)
		if err != nil {
			restErr := resterrors.ErrInvalidQueryParameter.New()
			restErr.AddDetailInfo(err.Error())
			c.JSON(rw, req, 400, c.getResponseBodyErr(400, restErr))
			return
		}

		var params *repository.ListParameters
		if c.listPolicy != nil {
			params = &repository.ListParameters{}
			err = forms.BindQuery(req, params, &c.listPolicy.Policy)
			if err != nil {
				restErr := resterrors.ErrInvalidQueryParameter.New()
				restErr.AddDetailInfo(err.Error())
				c.JSON(rw, req, 400, c.getResponseBodyErr(400, restErr))
				return
			}
		}
		var result interface{}
		var count int
		var dbErr *dberrors.Error

		if params != nil {
			if !params.ContainsParameters() {
				params.Limit = c.listPolicy.DefaultLimit
			}
			result, dbErr = c.repo.ListWithParams(obj, params)
		} else {
			result, dbErr = c.repo.List(obj)
		}
		if dbErr != nil {
			c.handleDBError(rw, req, dbErr)
			return
		}

		body := c.getResponseBodyContent(200, result)

		if c.listPolicy.WithCount {
			// Get Count for given collection
			count, dbErr = c.repo.Count(model)
			if dbErr != nil {
				c.handleDBError(rw, req, dbErr)
				return
			}
			type Count int
			// Add as 'count' to the body Content
			body.AddContent(Count(count))
		}

		c.JSON(rw, req, 200, body)
	}
}

func (c *GenericHandler) Update(model interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		obj := refutils.ObjOfPtrType(model)

		err := c.idSetFunc(req, obj)
		if err != nil {
			restErr := resterrors.ErrInternalError.New()
			c.JSON(rw, req, 500, c.getResponseBodyErr(500, restErr))
			return
		}

		err = forms.BindJSON(req, obj, c.jsonPolicy)
		if err != nil {
			restErr := resterrors.ErrInvalidJSONDocument.New()
			restErr.AddDetailInfo(err.Error())
			c.JSON(rw, req, 400, c.getResponseBodyErr(400, restErr))
		}

		dbErr := c.repo.Update(obj)
		if dbErr != nil {
			c.handleDBError(rw, req, dbErr)
			return
		}

		c.JSON(rw, req, 200, c.getResponseBodyContent(200, obj))
		return
	}
}

func (c *GenericHandler) Patch(model interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		// whereObj is the object that defines field to query the model
		var status int

		whereObj := refutils.ObjOfPtrType(model)
		if err := c.idSetFunc(req, whereObj); err != nil {
			restErr := resterrors.ErrInternalError.New()
			status = 500
			c.JSON(rw, req, status, c.getResponseBodyErr(status, restErr))
			return
		}

		obj := refutils.ObjOfPtrType(model)

		if err := forms.BindJSON(req, obj, c.jsonPolicy); err != nil {
			restErr := resterrors.ErrInvalidJSONDocument.New()
			restErr.AddDetailInfo(err.Error())
			status = 400
			c.JSON(rw, req, status, c.getResponseBodyErr(status, restErr))
			return
		}

		dbErr := c.repo.Patch(obj, whereObj)
		if dbErr != nil {
			c.handleDBError(rw, req, dbErr)
			return
		}

		result, dbErr := c.repo.Get(whereObj)
		if dbErr != nil {
			c.handleDBError(rw, req, dbErr)
			return
		}

		status = 200
		c.JSON(rw, req, status, c.getResponseBodyContent(status, result))
	}
}

func (c *GenericHandler) Delete(model interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		var status int

		whereObj := refutils.ObjOfPtrType(model)

		if err := c.idSetFunc(req, whereObj); err != nil {
			restErr := resterrors.ErrInternalError.New()
			status = 500
			c.JSON(rw, req, status, c.getResponseBodyErr(status, restErr))
			return
		}

		obj := refutils.ObjOfPtrType(model)
		dbErr := c.repo.Delete(obj, whereObj)
		if dbErr != nil {
			c.handleDBError(rw, req, dbErr)
			return
		}

		status = 200
		c.JSON(rw, req, status, c.getResponseBodyContent(204))
	}
}

func (c *GenericHandler) JSON(
	rw http.ResponseWriter,
	req *http.Request,
	status int,
	body response.Responser,
) {
	if status == 0 {
		status = 200
	}
	marshaledBody, err := json.Marshal(body)
	if err != nil {
		body = (&response.DefaultBody{}).NewErrored().WithErrors(resterrors.ErrInternalError.New())
		status = 500
		marshaledBody, _ = json.Marshal(body)
	}
	rw.WriteHeader(status)
	rw.Write(marshaledBody)
	return
}

func (c *GenericHandler) handleDBError(
	rw http.ResponseWriter,
	req *http.Request,
	dbError *dberrors.Error,
) {
	var isInternal bool
	var status int
	restErr, err := c.errHandler.Handle(dbError)
	if err != nil {
		isInternal = true
		restErr = resterrors.ErrInternalError.New()
	} else {
		isInternal = restErr.Compare(resterrors.ErrInternalError)
	}

	if isInternal {
		status = 500
	} else {
		status = 400
	}
	c.JSON(rw, req, status, c.getResponseBodyErr(status, restErr))
	return
}

func (c *GenericHandler) getResponseBodyErr(
	status int, errs ...*resterrors.Error,
) response.Responser {
	body := c.responseBody.NewErrored().WithErrors(errs...)
	if body, ok := body.(response.StatusResponser); ok {
		body.WithStatus(status)
	}
	return body
}

func (c *GenericHandler) getResponseBodyContent(
	status int, content ...interface{},
) response.Responser {
	body := c.responseBody.New().WithContent(content...)
	if body, ok := body.(response.StatusResponser); ok {
		body.WithStatus(status)
	}
	return body
}
