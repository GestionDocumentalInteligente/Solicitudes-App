package sdkhcl

import (
	"github.com/spf13/viper"

	"github.com/teamcubation/sg-backend/pkg/rest/clients/net-http/defs"
)

func Bootstrap(tokenEndPoint, clientID, clientSecret string, additionalParams map[string]string) (defs.Client, error) {
	if tokenEndPoint == "" {
		tokenEndPoint = viper.GetString("HTTP_CLIENT_ENDPOINT_KEY")
	}
	if clientID == "" {
		tokenEndPoint = viper.GetString("HTTP_CLIENT_CLIENT_ID")
	}
	if clientSecret == "" {
		clientSecret = viper.GetString("HTTP_CLIENT_SECRET")
	}
	if additionalParams == nil {
		additionalParams = viper.GetStringMapString("HTTP_CLIENT_ADD_PARAMS")
	}

	config := newConfig(
		tokenEndPoint,
		clientID,
		clientSecret,
		additionalParams,
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newClient(config)
}
