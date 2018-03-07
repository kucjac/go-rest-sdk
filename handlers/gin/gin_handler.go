package ginhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/inflection"
	"github.com/kucjac/go-rest-sdk"
	"github.com/kucjac/go-rest-sdk/errors"
	"log"
)

type GinJsonHandler struct {
	repository restsdk.GenericRepository
}

func (g *GinJsonHandler) Create(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		g.checkRepository()

		obj := restsdk.ObjOfPtrType(model)
		err := restsdk.BindJSON(c.Request, obj, &restsdk.FormPolicy{FailOnError: true})
		if err != nil {
			resErr := &resterrors.ErrInvalidJSONRequest
			resErr.ExtendDetail(err.Error())
			restsdk.ResponseWithError(400, resErr)
			return
		}

		// Create using provided repository
		err = g.repository.Create(obj)
		if err != nil {
			if restErr, ok := err.(*resterrors.ResponseError); ok {

				//If the given error is responseError check wether it is caused on client-side
				_, IsClientError := resterrors.ClientErrorCodes[restErr.Code]
				if IsClientError {

					//If this is client side error send it to response
					c.JSON(400, restsdk.ResponseWithError(400, restErr))
					return
				}
			}

			// If the error is not of type ResponseError or is not client side errors send internal server error
			// as a response
			c.JSON(500, restsdk.ResponseWithError(500, &resterrors.ErrInternalServerError))
			return
		}
		response := restsdk.ResponseWithOk()
		response.AddResult(restsdk.StructName(obj), obj)
		c.JSON(200, response)
	}
}

func (g *GinJsonHandler) Get(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		g.checkRepository()
		modelName := restsdk.StructName(model)
		modelID := c.Param(modelName + "id")

		if modelID == "" {
			c.JSON(400, restsdk.ResponseWithError(400, &resterrors.ErrBadRequestNoID))
			return
		}

		obj := restsdk.ObjOfPtrType(model)

		result, err := g.repository.Get(obj)
		if err != nil {
			// Check if error is of known type
			if rErr, ok := err.(*resterrors.ResponseError); ok {
				//Check if it is client side error
				_, clientError := resterrors.ClientErrorCodes[rErr.Code]
				if clientError {
					c.JSON(400, restsdk.ResponseWithError(400, rErr))
					return
				}
			}
			// Internal Server Error in all other types
			c.JSON(500, restsdk.ResponseWithError(500, &resterrors.ErrInternalServerError))
			return
		}

		response := restsdk.ResponseWithOk()
		response.AddResult(modelName, result)
		c.JSON(200, response)
		return
	}
}

func (g *GinJsonHandler) List(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestObject := restsdk.ObjOfType(model)
		err := restsdk.BindQuery(c.Request, &requestObject, nil)
		if err != nil {
			c.JSON(400, restsdk.ResponseWithError(400, &resterrors.ErrInvalidQueryParameters))
			return
		}

		meta := restsdk.ListParameters{}
		err = restsdk.BindQuery(c.Request, &meta, nil)
		if err != nil {
			c.JSON(400, restsdk.ResponseWithError(400, &resterrors.ErrInvalidQueryParameters))
			return
		}
		var result interface{}
		if meta.ContainsParameters() {
			result, err = g.repository.ListWithParams(requestObject, &meta)
		} else {
			result, err = g.repository.List(requestObject)
		}
		if err != nil {
			if rErr, ok := err.(*resterrors.ResponseError); ok {
				// Check client Error
				_, clientError := resterrors.ClientErrorCodes[rErr.Code]
				if clientError {
					c.JSON(400, restsdk.ResponseWithError(400, rErr))
					return
				}
			}
			c.JSON(500, restsdk.ResponseWithError(500, &resterrors.ErrInternalServerError))
			return
		}

		response := restsdk.ResponseWithOk()
		response.AddResult(inflection.Plural(restsdk.StructName(model)), &result)
		c.JSON(200, response)
	}
}

func (g *GinJsonHandler) Update(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (g *GinJsonHandler) Patch(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (g *GinJsonHandler) Delete(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (g *GinJsonHandler) checkRepository() {
	if g.repository == nil {
		log.Fatal("No repository set for handler")
	}
}
