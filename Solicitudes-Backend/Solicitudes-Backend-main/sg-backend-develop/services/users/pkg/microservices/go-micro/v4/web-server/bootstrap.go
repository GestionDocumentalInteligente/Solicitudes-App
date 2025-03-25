package sdkgomicro

import (
	"github.com/spf13/viper"

	defs "github.com/teamcubation/sg-users/pkg/microservices/go-micro/v4/web-server/defs"
)

func Bootstrap(webRouter interface{}) (defs.Server, error) {
	config := newConfig(
		webRouter,
		viper.GetString("WEB_SERVER_NAME"),
		viper.GetString("CONSUL_ADDRESS"),
		viper.GetString("WEB_SERVER_HOST"),
		viper.GetInt("WEB_SERVER_PORT"),
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newServer(config)
}
