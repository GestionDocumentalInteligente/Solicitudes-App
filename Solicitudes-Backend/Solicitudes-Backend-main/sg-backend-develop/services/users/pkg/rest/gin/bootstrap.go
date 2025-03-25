package sdkgin

import (
	"github.com/spf13/viper"

	defs "github.com/teamcubation/sg-users/pkg/rest/gin/defs"
)

func Bootstrap() (defs.Server, error) {
	config := newConfig(
		viper.GetString("WEB_SERVER_PORT"),
		viper.GetString("API_VERSION"),
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newServer(config)
}
