package ports

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server interface {
	RunServer() error
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
