package sdkgin

import (
	"os"

	defs "github.com/teamcubation/sg-mailing-api/pkg/rest/gin/defs"
)

func Bootstrap() (defs.Server, error) {
	config := newConfig(
		os.Getenv("WEB_SERVER_PORT"),
		os.Getenv("API_VERSION"),
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newServer(config)
}
