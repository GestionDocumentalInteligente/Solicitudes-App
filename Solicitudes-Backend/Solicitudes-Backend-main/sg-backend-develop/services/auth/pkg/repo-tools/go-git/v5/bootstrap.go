package sdkgogit

import (
	"os"

	ports "github.com/teamcubation/sg-auth/pkg/repo-tools/go-git/v5/ports"
)

func Bootstrap(repoRemoteUrlKey, repoLocalPathKey, repoBranchKey string) (ports.Client, error) {
	config := newConfig(
		os.Getenv(repoRemoteUrlKey),
		os.Getenv(repoLocalPathKey),
		os.Getenv(repoBranchKey),
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newClient(config)
}
