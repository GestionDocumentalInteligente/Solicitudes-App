package sdkhcl

import (
	"log"
	"net/url"

	"github.com/teamcubation/sg-backend/pkg/rest/clients/net-http/defs"
)

type config struct {
	tokenEndPoint    string
	clientID         string
	clientSecret     string
	additionalParams url.Values
}

func newConfig(tokenEndPoint, clientID, clientSecret string, additionalParams map[string]string) defs.Config {
	c := &config{
		tokenEndPoint:    tokenEndPoint,
		clientID:         clientID,
		clientSecret:     clientSecret,
		additionalParams: make(url.Values),
	}
	for key, value := range additionalParams {
		c.setAdditionalParam(key, value)
	}
	return c
}

func (c *config) GetTokenEndpoint() string {
	return c.tokenEndPoint
}

func (c *config) GetClientID() string {
	return c.clientID
}

func (c *config) GetClientSecret() string {
	return c.clientSecret
}

func (c *config) GetAdditionalParams() url.Values {
	return c.additionalParams
}

func (c *config) setAdditionalParam(key, value string) {
	c.additionalParams.Set(key, value)
}

func (c *config) Validate() error {
	if c.tokenEndPoint == "" {
		log.Println("Warning: token endpoint is not configured")
	}
	if c.clientID == "" {
		log.Println("Warning: client ID is not configured")
	}
	if c.clientSecret == "" {
		log.Println("Warning: secret is not configured")
	}
	return nil
}
