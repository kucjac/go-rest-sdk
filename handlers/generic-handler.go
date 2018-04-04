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
	"log"
	"net/http"
	"os"
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
	// Repository for given handler
	Repo repository.Repository

	ErrHandler *errhandler.ErrorHandler
	// ResponseBody - body used by in responses
	ResponseBody response.Responser

	// logger
	Log logger.ExtendedLeveledLogger

	// QueryPolicy - current policy for binding queries using BindQuery
	QueryPolicy *forms.BindPolicy

	// ParamPolicy - policy used for binding Parameters with BindParams
	ParamPolicy *forms.ParamPolicy
	// GetParams - ParamGetterFunction used for getting parameters used by third-party routers
	GetParams forms.ParamGetterFunc
	// with params specify if given route should bind parameters
	UseURLParams bool

	//ListParams
	ListParams *repository.ListParameters
	//UseCount flag for List method - defines if the response should include count of given
	//collection
	IncludeListCount bool
}

type SetIDFunc func(req *http.Request, model interface{}) error

// New creates GenericHandler for given
func New(repo repository.Repository,
	errHandler *errhandler.ErrorHandler,
	responseBody response.Responser,
	logs logger.ExtendedLeveledLogger,
) (handler *GenericHandler, err error) {
	if repo == nil || errHandler == nil {
		return nil, errors.New("Repository and errorHandler cannot be nil.")
	}
	if responseBody == nil {
		responseBody = &response.DefaultBody{}
	}
	handler = &GenericHandler{
		Repo:         repo,
		ErrHandler:   errHandler,
		ResponseBody: responseBody,
	}
	if logs == nil {
		handler.Log, _ = logger.NewLoggerWrapper(logger.NewBasicLogger(os.Stderr, "", log.Ldate))
	}
	return handler, nil
}

// New creates a copy of given handler and returns it.
func (c *GenericHandler) New() *GenericHandler {
	handlerCopy := *c
	return &handlerCopy
}

// WithQueryPolicy sets the query policy for given handler.
// Returns given handler so it can be used in a callback manner
func (c *GenericHandler) WithQueryPolicy(policy *forms.BindPolicy) *GenericHandler {
	c.QueryPolicy = policy
	return c
}

// WithParamPolicy sets the param policy for given handler.
// Returns given handler so it can be used in a callback manner
func (c *GenericHandler) WithParamPolicy(policy *forms.ParamPolicy) *GenericHandler {
	c.ParamPolicy = policy
	return c
}

// WithListParameters sets the ListParameters for the Select method
// Returns given handler so it can be used in a callback manner
func (c *GenericHandler) WithListParameters(
	params *repository.ListParameters,
) *GenericHandler {
	c.ListParams = params
	return c
}

// WithSelectCount sets the IncludeListCount flag for the handler
func (c *GenericHandler) WithSelectCount(
	includeCount bool,
) *GenericHandler {
	c.IncludeListCount = includeCount
	return c
}

// WithParams sets the given handler so that is binds the routing parameters to the model.
// Returns given handler so it can be used in a callback manner
func (c *GenericHandler) WithURLParams(
	useParams bool,
) *GenericHandler {
	c.UseURLParams = useParams
	return c
}

//WithParamGetterFunc sets the param getter func for given handler
func (c *GenericHandler) WithParamGetterFunc(
	GetParams forms.ParamGetterFunc,
) *GenericHandler {
	c.GetParams = GetParams
	return c
}

// WithResponseBody sets the response body for the GenericHandler
func (c *GenericHandler) WithResponseBody(
	body response.Responser,
) *GenericHandler {
	c.ResponseBody = body
	return c
}

// Create is a chiHandler HandlerFunc for creating new restful model records
// if the flag 'WithParams' is set to true and no GetParams is set for handler
// the handler will panic
func (c *GenericHandler) Create(model interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		var status int
		obj := refutils.ObjOfPtrType(model)

		err := forms.BindJSON(req, obj)
		if err != nil {
			restErr := resterrors.ErrInvalidJSONDocument.New()
			restErr.AddDetailInfo(err.Error())
			status = 400
			c.JSON(rw, req, status, c.getResponseBodyErr(status, restErr))
			return
		}

		// Set parameter if WithParams flag is set to true
		if c.UseURLParams {
			// bind params
			err := forms.BindParams(req, obj, c.GetParams, c.ParamPolicy)
			// if error occured - either the policy FailOnError is set or cannot set
			// other parameters
			if err != nil {
				c.Log.Errorf("%v: %s", req.URL.Path, err)
				restErr := resterrors.ErrInternalError.New()
				status = 500
				c.JSON(rw, req, status, c.getResponseBodyErr(status, restErr))
				return
			}
		}

		dbErr := c.Repo.Create(obj)
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

		if c.UseURLParams {
			err := forms.BindParams(req, obj, c.GetParams, c.ParamPolicy)
			if err != nil {
				restErr := resterrors.ErrInternalError.New()
				c.Log.Errorf("%v: %v", req.URL.Path, err)
				c.JSON(rw, req, 500, c.getResponseBodyErr(500, restErr))
				return
			}
		}

		result, dbErr := c.Repo.Get(obj)
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

		// Bind Query
		err := forms.BindQuery(req, obj, c.QueryPolicy)
		if err != nil {
			restErr := resterrors.ErrInvalidQueryParameter.New()
			restErr.AddDetailInfo(err.Error())
			c.JSON(rw, req, 400, c.getResponseBodyErr(400, restErr))
			return
		}

		var params *repository.ListParameters
		// Set List Parameters
		if c.ListParams != nil {
			params = new(repository.ListParameters)
			err = forms.BindQuery(req, params, c.QueryPolicy)
			if err != nil {
				restErr := resterrors.ErrInvalidQueryParameter.New()
				restErr.AddDetailInfo(err.Error())
				c.JSON(rw, req, 400, c.getResponseBodyErr(400, restErr))
				return
			}
		}

		// set URL parameters
		if c.UseURLParams {
			err := forms.BindParams(req, obj, c.GetParams, c.ParamPolicy)
			if err != nil {
				restErr := resterrors.ErrInternalError.New()
				c.Log.Errorf("%v: %v", req.URL.Path, err)
				c.JSON(rw, req, 500, c.getResponseBodyErr(500, restErr))
				return
			}
		}

		var result interface{}
		var dbErr *dberrors.Error

		if params != nil {
			if !params.ContainsParameters() {
				params.Limit = c.ListParams.Limit
			}
			result, dbErr = c.Repo.ListWithParams(obj, params)
		} else {
			result, dbErr = c.Repo.List(obj)
		}
		if dbErr != nil {
			c.handleDBError(rw, req, dbErr)
			return
		}

		body := c.getResponseBodyContent(200, result)

		// CollectionCount
		var collectionCount int
		if c.IncludeListCount {
			// Get Count for given collection
			collectionCount, dbErr = c.Repo.Count(model)
			if dbErr != nil {
				c.handleDBError(rw, req, dbErr)
				return
			}
			// Add as 'count' to the body Content
			type Count int

			body.AddContent(Count(collectionCount))
		}

		c.JSON(rw, req, 200, body)
	}
}

func (c *GenericHandler) Update(model interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		obj := refutils.ObjOfPtrType(model)

		err := forms.BindJSON(req, obj)
		if err != nil {
			restErr := resterrors.ErrInvalidJSONDocument.New()
			restErr.AddDetailInfo(err.Error())
			c.JSON(rw, req, 400, c.getResponseBodyErr(400, restErr))
			return
		}

		// set URL parameters
		if c.UseURLParams {
			err := forms.BindParams(req, obj, c.GetParams, c.ParamPolicy)
			if err != nil {
				restErr := resterrors.ErrInternalError.New()
				c.Log.Errorf("%v: %v", req.URL.Path, err)
				c.JSON(rw, req, 500, c.getResponseBodyErr(500, restErr))
				return
			}
		}

		dbErr := c.Repo.Update(obj)
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

		// set URL parameters
		if c.UseURLParams {
			err := forms.BindParams(req, whereObj, c.GetParams, c.ParamPolicy)
			if err != nil {
				restErr := resterrors.ErrInternalError.New()
				c.Log.Errorf("%v: %v", req.URL.Path, err)
				c.JSON(rw, req, 500, c.getResponseBodyErr(500, restErr))
				return
			}
		}

		obj := refutils.ObjOfPtrType(model)

		if err := forms.BindJSON(req, obj); err != nil {
			restErr := resterrors.ErrInvalidJSONDocument.New()
			restErr.AddDetailInfo(err.Error())
			status = 400
			c.JSON(rw, req, status, c.getResponseBodyErr(status, restErr))
			return
		}

		dbErr := c.Repo.Patch(obj, whereObj)
		if dbErr != nil {
			c.handleDBError(rw, req, dbErr)
			return
		}

		result, dbErr := c.Repo.Get(whereObj)
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

		// set URL parameters
		if c.UseURLParams {
			err := forms.BindParams(req, whereObj, c.GetParams, c.ParamPolicy)
			if err != nil {
				restErr := resterrors.ErrInternalError.New()
				c.Log.Errorf("%v: %v", req.URL.Path, err)
				c.JSON(rw, req, 500, c.getResponseBodyErr(500, restErr))
				return
			}
		}

		obj := refutils.ObjOfPtrType(model)
		dbErr := c.Repo.Delete(obj, whereObj)
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
		c.Log.Errorf("On: %s route, an error occurred while marshaling: %v", req.URL.Path, err)
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
	restErr, err := c.ErrHandler.Handle(dbError)
	if err != nil {
		isInternal = true
		c.Log.Errorf("%v: %v", req.URL.Path, err)
		restErr = resterrors.ErrInternalError.New()
	} else {
		isInternal = restErr.Compare(resterrors.ErrInternalError)
	}

	if isInternal {
		c.Log.Errorf("%v: %v", req.URL.Path, err)
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
	body := c.ResponseBody.NewErrored().WithErrors(errs...)
	if body, ok := body.(response.StatusResponser); ok {
		body.WithStatus(status)
	}
	return body
}

func (c *GenericHandler) getResponseBodyContent(
	status int, content ...interface{},
) response.Responser {
	body := c.ResponseBody.New().WithContent(content...)
	if body, ok := body.(response.StatusResponser); ok {
		body.WithStatus(status)
	}
	return body
}
