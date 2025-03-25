package sdkgin

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	ports "github.com/teamcubation/sg-auth/pkg/rest/gin/ports"
)

var (
	instance  ports.Server
	once      sync.Once
	initError error
)

type server struct {
	router *gin.Engine
	config ports.Config
}

func newServer(config ports.Config) (ports.Server, error) {
	once.Do(func() {
		err := config.Validate()
		if err != nil {
			initError = err
			return
		}

		gin.SetMode(gin.ReleaseMode)
		r := gin.Default()
		instance = &server{
			config: config,
			router: r,
		}
	})
	return instance, initError
}

func (server *server) RunServer() error {
	return server.router.Run(":" + server.config.GetRouterPort())
}

func (server *server) GetRouter() *gin.Engine {
	return server.router
}

func (server *server) GetApiVersion() string {
	return server.config.GetApiVersion()
}

// WrapH envuelve un http.Handler en un gin.HandlerFunc.
func (server *server) WrapH(h http.Handler) gin.HandlerFunc {
	return gin.WrapH(h)
}
