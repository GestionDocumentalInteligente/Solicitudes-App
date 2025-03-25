package defs

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server interface {
	RunServer(context.Context) error
	GetApiVersion() string
	GetRouter() *gin.Engine
	WrapH(h http.Handler) gin.HandlerFunc
}

type Config interface {
	GetRouterPort() string
	SetRouterPort(string)
	GetApiVersion() string
	SetApiVersion(string)
	Validate() error
}
