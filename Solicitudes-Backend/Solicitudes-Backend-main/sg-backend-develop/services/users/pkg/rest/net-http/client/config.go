package sdkhclnt

import (
	"fmt"
	"net/url"

	"github.com/teamcubation/sg-users/pkg/rest/net-http/client/defs"
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
		c.SetAdditionalParam(key, value)
	}
	return c
}

func (c *config) GetTokenEndpoint() string {
	return c.tokenEndPoint
}

func (c *config) SetTokenEndpoint(endpoint string) {
	c.tokenEndPoint = endpoint
}

func (c *config) GetClientID() string {
	return c.clientID
}

func (c *config) SetClientID(id string) {
	c.clientID = id
}

func (c *config) GetClientSecret() string {
	return c.clientSecret
}

func (c *config) SetClientSecret(secret string) {
	c.clientSecret = secret
}

func (c *config) GetAdditionalParams() url.Values {
	return c.additionalParams
}

func (c *config) SetAdditionalParams(params url.Values) {
	c.additionalParams = params
}

func (c *config) SetAdditionalParam(key, value string) {
	c.additionalParams.Set(key, value)
}

func (c *config) Validate() error {
	if c.tokenEndPoint == "" {
		return fmt.Errorf("token endpoint is not configured")
	}
	if c.clientSecret == "" {
		return fmt.Errorf("client secret is not configured")
	}
	return nil
}
