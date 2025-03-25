package sdkaws

import (
	"fmt"

	"github.com/teamcubation/sg-backend/pkg/aws/localstack/defs"
)

type Config struct {
	AwsAccessKeyID     string
	AwsSecretAccess    string
	AwsRegion          string
	LocalStackEndpoint string
}

func newConfig(awsAccessKeyID, awsSecretAccessKey, awsRegion, localStackEndpoint string) defs.Config {
	return &Config{
		AwsAccessKeyID:     awsAccessKeyID,
		AwsSecretAccess:    awsSecretAccessKey,
		AwsRegion:          awsRegion,
		LocalStackEndpoint: localStackEndpoint,
	}
}

func (c *Config) Validate() error {
	if c.AwsAccessKeyID == "" {
		return fmt.Errorf("AWS_ACCESS_KEY_ID is required")
	}
	if c.AwsSecretAccess == "" {
		return fmt.Errorf("AWS_SECRET_ACCESS_KEY is required")
	}
	if c.AwsRegion == "" {
		return fmt.Errorf("AWS_REGION is required")
	}
	if c.LocalStackEndpoint == "" {
		return fmt.Errorf("LOCALSTACK_ENDPOINT is required")
	}
	return nil
}

func (c *Config) GetAwsAccessKeyID() string {
	return c.AwsAccessKeyID
}

func (c *Config) GetAwsSecretAccessKey() string {
	return c.AwsSecretAccess
}

func (c *Config) GetAwsRegion() string {
	return c.AwsRegion
}

func (c *Config) GetLocalStackEndpoint() string {
	return c.LocalStackEndpoint
}
