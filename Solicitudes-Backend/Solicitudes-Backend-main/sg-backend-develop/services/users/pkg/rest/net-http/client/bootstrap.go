package sdkhclnt

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/teamcubation/sg-users/pkg/rest/net-http/client/defs"
)

func Bootstrap(tokenEndPointKey, clientIDKey, clientSecretKey, addParamsKey string) (defs.Client, defs.Config, error) {
	tokenEndPoint := viper.GetString(tokenEndPointKey)
	if tokenEndPoint == "" {
		return nil, nil, fmt.Errorf("token endpoint is empty. Check if %s environment variable is set", tokenEndPointKey)
	}

	config := newConfig(
		tokenEndPoint,
		viper.GetString(clientIDKey),
		viper.GetString(clientSecretKey),
		viper.GetStringMapString(addParamsKey),
	)

	if err := config.Validate(); err != nil {
		return nil, nil, err
	}

	client, err := newClient(config)
	if err != nil {
		return nil, nil, err
	}

	return client, config, nil
}
