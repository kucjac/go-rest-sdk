package handlers

import (
	"encoding/json"
	"errors"
	"github.com/kucjac/go-rest-sdk/dberrors"
	"github.com/kucjac/go-rest-sdk/errhandler"
	"github.com/kucjac/go-rest-sdk/forms"
	"github.com/kucjac/go-rest-sdk/logger"
	"github.com/kucjac/go-rest-sdk/refutils"
	"github.com/kucjac/go-rest-sdk/repository"
	"github.com/kucjac/go-rest-sdk/response"
	"github.com/kucjac/go-rest-sdk/resterrors"
	"net/http"
)

var (
	ErrIncorrectModel         = errors.New("Incorrect model route path provided.")
	ErrIncorrectCustomContext = errors.New("Incorrect custom context type.")
	ErrNoParamGetterFuncSet   = errors.New("No ParamGetterFunc set.")
)

// GenericHandler is a structure that is used to build basic
// CRUD operations for RESTful API's.
//

type GenericHandler struct {
	// repository for given handler
	repo repository.Repository

	errHandler *errhandler.ErrorHandler

	// logger
	log logger.GenericLogger

	// queryPolicy - current policy for binding queries using BindQuery
	queryPolicy *forms.Policy

	// jsonPolicy - policy for binding json using BindJSON
	jsonPolicy *forms.Policy

	// listPolicy - policy used or bindings Lists
	listPolicy *forms.ListPolicy

	// paramPolicy - policy used for binding Parameters with BindParams
	paramPolicy *forms.ParamPolicy

	// responseBody - body used by in responses
	responseBody response.Responser

	// getParam - ParamGetterFunction used for getting parameters used by third-party routers
	getParam forms.ParamGetterFunc

	// with params specify if given route should bind parameters
	withParams bool
}

type SetIDFunc func(req *http.Request, model interface{}) error

// New creates GenericHandler for given
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

// New creates a copy of given handler and returns it.
func (c *GenericHandler) New() *GenericHandler {
	h := *c
	return &h
}

// WithQueryPolicy sets the query policy for given handler.
// Returns given handler so it can be used in a callback manner
func (c *GenericHandler) WithQueryPolicy(policy *forms.Policy) *GenericHandler {
	c.queryPolicy = policy
	return c
}

// WithJSONPolicy sets the policy used for BindJSON function.
// Returns given handler so it can be used in a callback manner
func (c *GenericHandler) WithJSONPolicy(policy *forms.Policy) *GenericHandler {
	c.jsonPolicy = policy
	return c
}

// WithListPolicy sets the police used in 'List' handler
// Returns given handler so it can be used in a callback manner
func (c *GenericHandler) WithListPolicy(policy *forms.ListPolicy) *GenericHandler {
	c.listPolicy = policy
	return c
}

// WithParamPolicy sets the param policy for given handler.
// Returns given handler so it can be used in a callback manner
func (c *GenericHandler) WithParamPolicy(policy *forms.ParamPolicy) *GenericHandler {
	c.paramPolicy = policy
	return c
}

// WithParams sets the given handler so that is binds the routing parameters to the model.
// Returns given handler so it can be used in a callback manner
func (c *GenericHandler) WithParams(useParams bool) *GenericHandler {
	c.withParams = useParams
	return c
}

// WithParamGetterFunc sets the ParamGetterFunc for given handler.
// Returns given handler so it can be used in a callback manner
func (c *GenericHandler) WithParamGetterFunc(getParam forms.ParamGetterFunc) *GenericHandler {
	c.getParam = getParam
	return c
}

func (c *GenericHandler) WithSetIDFunc(
	customIDFunc SetIDFunc,
) *GenericHandler {
	c.idSetFunc = customIDFunc
	return c
}

// Create is a chiHandler HandlerFunc for creating new restful model records
// if the flag 'withParams' is set to true and no getParam is set for handler
// the handler will panic
func (c *GenericHandler) Create(model interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		var status int
		obj := refutils.ObjOfPtrType(model)

		// Set parameter if withParams flag is set to true
		if c.withParams {

			// if no getParam function set
			if c.getParam == nil {
				c.log.Errorf("For %v, an error occurred: %v",
					req.URL.Path, ErrNoParamGetterFuncSet)
				restErr := resterrors.ErrInternalError.New()
				status = 500
				c.JSON(rw, req, status, c.getResponseBodyErr(status, restErr))
				return
			}

			// bind params
			err := forms.BindParams(req, model, c.getParam, c.paramPolicy)
			// if error occured - either the policy FailOnError is set or cannot set
			// other parameters
			if err != nil {
				restErr := resterrors.ErrInternalError.New()
				status = 500
				c.JSON(rw, req, status, c.getResponseBodyErr(status, restErr))
				return
			}
		}
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

		forms.BindParams(req, model, getParam, policy)
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
		c.log.Errorf("On: %v route, an error occurred: %v", req.URL.Path, err)
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
		c.log.Errorf("Error while handling DBError: %v", err)
		restErr = resterrors.ErrInternalError.New()
	} else {
		isInternal = restErr.Compare(resterrors.ErrInternalError)
	}

	if isInternal {
		c.log.Errorf("On the route Database error: %v", dbError)
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
