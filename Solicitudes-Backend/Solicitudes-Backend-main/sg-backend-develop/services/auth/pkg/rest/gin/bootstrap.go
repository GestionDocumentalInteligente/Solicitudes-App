package sdkgin

import (
	"os"

	ports "github.com/teamcubation/sg-auth/pkg/rest/gin/ports"
)

func Bootstrap() (ports.Server, error) {
	config := newConfig(
		os.Getenv("WEB_SERVER_PORT"),
		os.Getenv("API_VERSION"),
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newServer(config)
}
