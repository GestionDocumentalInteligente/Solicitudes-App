package sdkcnfldr

import (
	"errors"
	"fmt"

	sdkgodotenv "github.com/teamcubation/sg-users/pkg/config/godotenv"
	sdkviper "github.com/teamcubation/sg-users/pkg/config/viper"
	sdktools "github.com/teamcubation/sg-users/pkg/tools"
)

func LoadConfig(filePaths ...string) error {
	if len(filePaths) == 0 {
		return errors.New("no file paths provided")
	}

	foundFiles, err := sdktools.FilesFinder(filePaths...)
	if err != nil {
		return fmt.Errorf("FilesFinder failed to find files: %w", err)
	}

	if err := sdkgodotenv.LoadConfig(foundFiles...); err != nil {
		return fmt.Errorf("sdkcnfldr: %v", err)
	}

	if err := sdkviper.LoadConfig(foundFiles...); err != nil {
		return fmt.Errorf("sdkcnfldr: %v", err)
	}

	return nil
}
