package ginhandler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kucjac/go-rest-sdk"
)

type GinJsonHandler struct {
	repository restsdk.GenericRepository
}

func (g *GinJsonHandler) Create(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		obj := &model
		err := restsdk.BindJSON(c.Request, &obj, &restsdk.FormPolicy{FailOnError: true})
		if err != nil {
			resErr := &restsdk.ErrInvalidJSONRequest
			resErr.Detail += fmt.Sprintf(" %s", err.Error())
			restsdk.ResponseWithError(400, resErr)
			return
		}
		err = g.repository.Create(obj)
		if err != nil {
			// handle errors
			restsdk.ResponseWithError(400, err)
			return
		}
		restsdk.ResponseWithOk()
	}
}

func (g *GinJsonHandler) Get(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (g *GinJsonHandler) List(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {

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
