package ginhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/kucjac/go-rest-sdk"
)

type GinJsonHandler struct {
	repository restsdk.GenericRepository
}

func (g *GinJsonHandler) Create(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := c.ShouldBindJSON(&model)
		if err != nil {

		}
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
