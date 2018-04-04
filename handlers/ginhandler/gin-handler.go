package ginhandler

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"

	"github.com/kucjac/go-rest-sdk/errhandler"
	"github.com/kucjac/go-rest-sdk/forms"
	"github.com/kucjac/go-rest-sdk/handlers"
	"github.com/kucjac/go-rest-sdk/logger"
	"github.com/kucjac/go-rest-sdk/repository"
	"github.com/kucjac/go-rest-sdk/response"

	"net/http"
)

// GinHandler the policies are set manually
type GinHandler struct {
	handlers.GenericHandler
}

// GinParamKey is a http.Request.Context key for gin.Params
// use it to save and get gin.Params from the context.
type GinParamKey struct{}

// GinParamGetterFunc is forms.ParamGetterFunc implementation for gin router.
// It gets the 'GinParamKey' from the request context.
// If not found an error would be returned.
func GinParamGetterFunc(param string, req *http.Request) (paramValue string, err error) {
	params, ok := req.Context().Value(GinParamKey{}).(gin.Params)
	if !ok {
		err = errors.New("No GinParamKey in context.")
		return
	}

	paramValue = params.ByName(param)
	return
}

// New Create a copy of GinHandler
func (g *GinHandler) New() *GinHandler {
	handler := *g
	return &handler
}

func (g *GinHandler) WithQueryPolicy(policy *forms.BindPolicy) *GinHandler {
	g.QueryPolicy = policy
	return g
}

func (g *GinHandler) WithParamPolicy(policy *forms.ParamPolicy) *GinHandler {
	g.ParamPolicy = policy
	return g
}

func (g *GinHandler) WithListParameters(
	params *repository.ListParameters,
	includeCount bool,
) *GinHandler {
	g.ListParams = params
	g.IncludeListCount = includeCount
	return g
}

func (g *GinHandler) WithURLParams(
	useParams bool,
) *GinHandler {
	g.UseURLParams = useParams
	return g
}

func (g *GinHandler) WithResponseBody(body response.Responser) *GinHandler {
	g.ResponseBody = body
	return g
}

// New creates GinHandler for given
func New(repo repository.Repository,
	errHandler *errhandler.ErrorHandler,
	responseBody response.Responser,
	Logger logger.ExtendedLeveledLogger,
) (ginHandler *GinHandler, err error) {
	var generic *handlers.GenericHandler
	generic, err = handlers.New(repo, errHandler, responseBody, Logger)
	if err != nil {
		return
	}
	generic.WithParamGetterFunc(GinParamGetterFunc)
	ginHandler = &GinHandler{
		GenericHandler: *generic,
	}
	return
}

// Create returns gin.handlerFunc that for given 'model' creates new entity
// on the base of the request json body.
func (g *GinHandler) Create(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		contextWithGinParams(c)
		g.GenericHandler.Create(model).ServeHTTP(c.Writer, c.Request)
	}
}

// Get is a JSON gin.HandlerFunc that gets given model entity
// with provided 'modelname_id' entity
// The model is taken from the repository based on its id and name
func (g *GinHandler) Get(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		contextWithGinParams(c)
		g.GenericHandler.Get(model).ServeHTTP(c.Writer, c.Request)
	}
}

func (g *GinHandler) List(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		contextWithGinParams(c)
		g.GenericHandler.List(model).ServeHTTP(c.Writer, c.Request)
	}
}

func (g *GinHandler) Update(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		contextWithGinParams(c)
		g.GenericHandler.Update(model).ServeHTTP(c.Writer, c.Request)
	}
}

func (g *GinHandler) Patch(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		contextWithGinParams(c)
		g.GenericHandler.Patch(model).ServeHTTP(c.Writer, c.Request)
	}
}

func (g *GinHandler) Delete(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		contextWithGinParams(c)
		g.GenericHandler.Patch(model).ServeHTTP(c.Writer, c.Request)
	}
}

func contextWithGinParams(c *gin.Context) {
	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, GinParamKey{}, c.Params)
	c.Request.WithContext(ctx)
}
