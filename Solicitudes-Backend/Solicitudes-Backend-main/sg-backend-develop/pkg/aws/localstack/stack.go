package sdkaws

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/teamcubation/sg-backend/pkg/aws/localstack/defs"
)

var (
	instance  defs.Stack
	once      sync.Once
	initError error
)

type stack struct {
	config    defs.Config
	awsConfig aws.Config
}

func newStack(c defs.Config) (defs.Stack, error) {
	once.Do(func() {
		svc := &stack{config: c}
		initError = svc.Connect()
		if initError != nil {
			instance = nil
		} else {
			instance = svc
		}
	})
	return instance, initError
}

func (s *stack) Connect() error {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(s.config.GetAwsRegion()),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			s.config.GetAwsAccessKeyID(), s.config.GetAwsSecretAccessKey(), "",
		)),
	)
	if err != nil {
		return err
	}
	s.awsConfig = awsConfig
	return nil
}

func (s *stack) GetConfig() aws.Config {
	return s.awsConfig
}
