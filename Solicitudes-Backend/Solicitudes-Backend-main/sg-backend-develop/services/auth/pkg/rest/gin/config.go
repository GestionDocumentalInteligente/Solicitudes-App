package sdkgin

import (
	"fmt"

	ports "github.com/teamcubation/sg-auth/pkg/rest/gin/ports"
)

type config struct {
	routerPort string
	ApiVersion string
}

func newConfig(routerPort, ApiVersion string) ports.Config {
	return &config{
		routerPort: routerPort,
		ApiVersion: ApiVersion,
	}
}

func (c *config) GetRouterPort() string {
	return c.routerPort
}

func (c *config) SetRouterPort(routerPort string) {
	c.routerPort = routerPort
}

func (c *config) GetApiVersion() string {
	return c.ApiVersion
}

func (c *config) SetApiVersion(ApiVersion string) {
	c.ApiVersion = ApiVersion
}

func (c *config) Validate() error {
	if c.routerPort == "" {
		return fmt.Errorf("router port is not configured")
	}
	return nil
}
