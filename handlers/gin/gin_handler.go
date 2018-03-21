package ginhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/inflection"
	"github.com/kucjac/go-rest-sdk/dberrors"
	"github.com/kucjac/go-rest-sdk/errhandler"
	"github.com/kucjac/go-rest-sdk/forms"
	"github.com/kucjac/go-rest-sdk/refutils"
	"github.com/kucjac/go-rest-sdk/repository"
	"github.com/kucjac/go-rest-sdk/response"
	"github.com/kucjac/go-rest-sdk/resterrors"
	"log"
)

// GinJSONHandler
type GinJSONHandler struct {
	repo       repository.GenericRepository
	errHandler *errhandler.ErrorHandler
}

func (g *GinJSONHandler) Create(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		obj := refutils.ObjOfPtrType(model)
		err := forms.BindJSON(c.Request, obj, &forms.FormPolicy{FailOnError: true})
		if err != nil {
			resErr := resterrors.ErrInvalidJSONDocument.New()
			resErr.AddDetailInfo(err.Error())
			c.JSON(400, response.NewWithError(400, resErr))
			return
		}

		// Create using provided repository
		dberr := g.repo.Create(obj)

		if dberr != nil {
			rErr, err := g.errHandler.Handle(dberr)
			if err != nil {
				c.JSON(500, response.NewWithError(500, resterrors.ErrInternalError.New()))
				return
			}
			isInternal := rErr.Compare(resterrors.ErrInternalError)
			if isInternal {
				c.JSON(500, response.NewWithError(500, rErr))
			} else {
				c.JSON(400, response.NewWithError(400, rErr))
			}
			return
		}

		body := response.New()
		body.AddContent(refutils.StructName(obj), obj)
		c.JSON(200, body)
	}
}

// Get is a JSON gin.HandlerFunc that gets given model entity
// with provided 'modelname_id' entity
// The model is taken from the repository based on its id and name
func (g *GinJSONHandler) Get(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get model name
		modelName := refutils.StructName(model)

		// modelID should be a url parameter as a ''
		modelID := c.Param(modelName + "_id")
		if modelID == "" {
			c.JSON(400, response.NewWithError(400, resterrors.ErrInvalidURI.New()))
			return
		}

		// create new object entity based on the model
		obj := refutils.ObjOfPtrType(model)

		// Set the model ID
		err := forms.SetID(obj, modelID)
		if err != nil {
			c.JSON(500, response.NewWithError(500, resterrors.ErrInternalError.New()))
			return
		}

		// get the specific model from the repository
		result, dberr := g.repo.Get(obj)
		if dberr != nil {
			var isInternal bool
			// Handle the error
			restErr, err := g.errHandler.Handle(dberr)
			if err != nil {
				isInternal = true
				restErr = resterrors.ErrInternalError.New()
			}
			// Internal Server Error in all other types
			if !isInternal {
				if isInternal = restErr.Compare(resterrors.ErrInternalError); !isInternal {
					c.JSON(400, response.NewWithError(400, restErr))
					return
				}
			}
			c.JSON(500, response.NewWithError(500, restErr))
			return
		}

		// Get body
		body := response.New()

		// Add content
		body.AddContent(modelName, result)

		// Marshal to json with http.Status - 200
		c.JSON(200, body)
		return
	}
}

func (g *GinJSONHandler) List(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create new request object for the list
		reqObj := refutils.ObjOfPtrType(model)

		// Bind URL Query to the req object
		err := forms.BindQuery(c.Request, reqObj, nil)
		if err != nil {
			c.JSON(400, response.NewWithError(400, resterrors.ErrInvalidQueryParameter.New()))
			return
		}

		// Bind URL Query to the list parameters
		meta := &repository.ListParameters{}
		err = forms.BindQuery(c.Request, meta, nil)
		if err != nil {
			c.JSON(400, response.NewWithError(400, resterrors.ErrInvalidQueryParameter.New()))
			return
		}

		var result interface{}
		var dberr *dberrors.Error

		if meta.ContainsParameters() {
			result, dberr = g.repo.ListWithParams(reqObj, meta)
		} else {
			result, dberr = g.repo.List(reqObj)
		}
		if dberr != nil {
			var isInternal bool
			restErr, err := g.errHandler.Handle(dberr)
			if err != nil {
				isInternal = true
				restErr = resterrors.ErrInternalError.New()
			}
			if !isInternal {
				if isInternal = restErr.Compare(resterrors.ErrInternalError); !isInternal {
					c.JSON(400, response.NewWithError(400, restErr))
					return
				}
			}
			c.JSON(500, response.NewWithError(500, restErr))
			return
		}

		body := response.New()
		body.AddContent(inflection.Plural(refutils.StructName(model)), result)
		c.JSON(200, body)
		return

	}
}

func (g *GinJSONHandler) Update(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (g *GinJSONHandler) Patch(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (g *GinJSONHandler) Delete(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (g *GinJSONHandler) checkRepository() {
	if g.repo == nil {
		log.Fatal("No repository set for handler")
	}
}
