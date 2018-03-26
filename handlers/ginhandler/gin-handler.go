package ginhandler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kucjac/go-rest-sdk/dberrors"
	"github.com/kucjac/go-rest-sdk/errhandler"
	"github.com/kucjac/go-rest-sdk/forms"
	"github.com/kucjac/go-rest-sdk/refutils"
	"github.com/kucjac/go-rest-sdk/repository"
	"github.com/kucjac/go-rest-sdk/response"
	"github.com/kucjac/go-rest-sdk/resterrors"
)

// JSONHandler the policies are set manually
type JSONHandler struct {
	repo             repository.Repository
	errHandler       *errhandler.ErrorHandler
	queryPolicy      *forms.Policy
	jsonPolicy       *forms.Policy
	parametersPolicy *forms.Policy
	responseBody     response.Responser
}

func (g *JSONHandler) New() *JSONHandler {
	return &(*g)
}

func (g *JSONHandler) WithQueryPolicy(policy *forms.Policy) *JSONHandler {
	g.queryPolicy = policy
	return g
}

func (g *JSONHandler) WithJSONPolicy(policy *forms.Policy) *JSONHandler {
	g.jsonPolicy = policy
	return g
}

func (g *JSONHandler) WithParamPolicy(policy *forms.Policy) *JSONHandler {
	g.parametersPolicy = policy
	return g
}

func (g *JSONHandler) WithResponseBody(body response.Responser) *JSONHandler {
	g.responseBody = body
	return g
}

// New creates JSONHandler for given
func New(repo repository.Repository,
	errHandler *errhandler.ErrorHandler,
	responseBody response.Responser,
) (*JSONHandler, error) {
	if repo == nil || errHandler == nil {
		return nil, errors.New("repository and errorHandler cannot be nil.")
	}
	if responseBody == nil {
		responseBody = &response.DefaultBody{}
	}
	ginHandler := &JSONHandler{
		repo:         repo,
		errHandler:   errHandler,
		responseBody: responseBody,
	}
	return ginHandler, nil
}

// Create returns gin.handlerFunc that for given 'model' creates new entity
// on the base of the request json body.
func (g *JSONHandler) Create(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {

		obj := refutils.ObjOfPtrType(model)
		err := forms.BindJSON(c.Request, obj, g.jsonPolicy)
		if err != nil {
			restErr := resterrors.ErrInvalidJSONDocument.New()
			restErr.AddDetailInfo(err.Error())
			c.JSON(400, g.getResponseBodyErr(400, restErr))
			return
		}

		// Create using provided repository
		dberr := g.repo.Create(obj)

		if dberr != nil {
			g.handleDBError(c, dberr)
			return
		}

		c.JSON(201, g.getResponseBodyCon(201, obj))
	}
}

// Get is a JSON gin.HandlerFunc that gets given model entity
// with provided 'modelname_id' entity
// The model is taken from the repository based on its id and name
func (g *JSONHandler) Get(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {

		// create new object entity based on the model
		obj := refutils.ObjOfPtrType(model)

		err := g.setModelID(c, obj)
		if err != nil {
			c.Error(err)
			restErr := resterrors.ErrInternalError.New()
			c.JSON(500, g.getResponseBodyErr(500, restErr))
			return
		}

		// get the specific model from the repository
		result, dberr := g.repo.Get(obj)
		if dberr != nil {
			g.handleDBError(c, dberr)
			return
		}

		// Marshal to json with http.Status - 200
		c.JSON(200, g.getResponseBodyCon(200, result))
		return
	}
}

func (g *JSONHandler) List(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create new request object for the list
		reqObj := refutils.ObjOfPtrType(model)

		// Bind URL Query to the req object
		err := forms.BindQuery(c.Request, reqObj, g.queryPolicy)
		if err != nil {
			restErr := resterrors.ErrInvalidQueryParameter.New()
			restErr.AddDetailInfo(err.Error())

			c.JSON(400, g.getResponseBodyErr(400, restErr))
			return
		}

		// Bind URL Query to the list parameters
		parameters := &repository.ListParameters{}
		err = forms.BindQuery(c.Request, parameters, g.parametersPolicy)
		if err != nil {
			restErr := resterrors.ErrInvalidQueryParameter.New()
			restErr.AddDetailInfo(err.Error())
			c.JSON(400, g.getResponseBodyErr(400, restErr))
			return
		}

		var result interface{}
		var count int
		var dberr *dberrors.Error

		if parameters.ContainsParameters() {
			result, dberr = g.repo.ListWithParams(reqObj, parameters)
		} else {
			result, dberr = g.repo.List(reqObj)
		}

		if dberr != nil {
			g.handleDBError(c, dberr)
			return
		}

		type Count int
		count, dberr = g.repo.Count(model)
		if dberr != nil {
			g.handleDBError(c, dberr)
			return
		}

		c.JSON(200, g.getResponseBodyCon(200, result, Count(count)))
		return
	}
}

func (g *JSONHandler) Update(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {

		reqObj := refutils.ObjOfPtrType(model)

		// Set model ID for given 'reqObj'
		if err := g.setModelID(c, reqObj); err != nil {
			c.Error(err)
			restErr := resterrors.ErrInternalError.New()
			c.JSON(500, g.responseBody.NewErrored().WithErrors(restErr))
			return
		}

		// BindJSON from the request
		err := forms.BindJSON(c.Request, reqObj, g.jsonPolicy)
		if err != nil {
			restErr := resterrors.ErrInvalidJSONDocument.New()
			restErr.AddDetailInfo(err.Error())

			c.JSON(400, g.getResponseBodyErr(400, restErr))
			return
		}

		dbErr := g.repo.Update(reqObj)
		if dbErr != nil {
			g.handleDBError(c, dbErr)
			return
		}

		// Response with the given requested object
		c.JSON(200, g.getResponseBodyCon(200, reqObj))
		return
	}
}

func (g *JSONHandler) Patch(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// whereObj is the object that defines fields that the model defines
		whereObj := refutils.ObjOfPtrType(model)
		if err := g.setModelID(c, whereObj); err != nil {
			c.Error(err)
			restErr := resterrors.ErrInternalError.New()
			c.JSON(500, g.getResponseBodyErr(500, restErr))
			return
		}

		reqObj := refutils.ObjOfPtrType(model)

		// BindJSON from the request
		err := forms.BindJSON(c.Request, reqObj, g.jsonPolicy)
		if err != nil {
			restErr := resterrors.ErrInvalidJSONDocument.New()
			restErr.AddDetailInfo(err.Error())

			c.JSON(400, g.getResponseBodyErr(400, restErr))
			return
		}

		dbErr := g.repo.Patch(reqObj, whereObj)
		if dbErr != nil {
			g.handleDBError(c, dbErr)
			return
		}

		// Response with the given requested object
		c.JSON(200, g.responseBody.New().WithContent(reqObj))
		return
	}
}

func (g *JSONHandler) Delete(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// where obj defines the query for deleting an object
		whereObj := refutils.ObjOfPtrType(model)

		// set the model id to 'whereObj'
		if err := g.setModelID(c, whereObj); err != nil {
			c.Error(err)
			restErr := resterrors.ErrInternalError.New()
			c.JSON(500, g.getResponseBodyErr(500, restErr))
			return
		}

		reqObj := refutils.ObjOfPtrType(model)

		dbErr := g.repo.Delete(reqObj, whereObj)
		if dbErr != nil {
			g.handleDBError(c, dbErr)
			return
		}

		// Response with the given requested object
		c.JSON(200, g.getResponseBodyCon(204))
		return
	}
}

func (g *JSONHandler) handleDBError(
	c *gin.Context,
	dbError *dberrors.Error,
) {
	if dbError != nil {
		var isInternal bool
		restErr, err := g.errHandler.Handle(dbError)
		if err != nil {
			c.Error(err)
			isInternal = true
			restErr = resterrors.ErrInternalError.New()
		} else {
			isInternal = restErr.Compare(resterrors.ErrInternalError)
		}
		if isInternal {
			c.JSON(500, g.getResponseBodyErr(500, restErr))
		} else {
			c.JSON(400, g.getResponseBodyErr(400, restErr))
		}
		return
	}
}

func (g *JSONHandler) setModelID(
	c *gin.Context,
	model interface{},
) (err error) {
	// Get model name
	modelName := refutils.ModelName(model)

	// modelID should be a url parameter as a ''
	modelID := c.Param(modelName)
	if modelID == "" {
		err = errors.New("Incorrect model route path provided.")
		return err
	}

	// SetID for given model
	err = forms.SetID(model, modelID)
	if err != nil {
		return err
	}
	return nil
}

func (g *JSONHandler) getResponseBodyErr(
	status int, errs ...*resterrors.Error,
) response.Responser {
	body := g.responseBody.NewErrored().WithErrors(errs...)
	if body, ok := body.(response.StatusResponser); ok {
		body.WithStatus(status)
	}
	return body
}

func (g *JSONHandler) getResponseBodyCon(
	status int, content ...interface{},
) response.Responser {
	body := g.responseBody.New().WithContent(content...)
	if body, ok := body.(response.StatusResponser); ok {
		body.WithStatus(status)
	}
	return body
}
