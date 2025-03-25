package defs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type Stack interface {
	Connect() error
	GetConfig() aws.Config
	NewSQSClient() SqsClient
}

type Config interface {
	GetAwsAccessKeyID() string
	GetAwsRegion() string
	GetAwsSecretAccessKey() string
	GetLocalStackEndpoint() string
	Validate() error
}

type SqsClient interface {
	GetOrCreateQueueURL(context.Context, string) (string, error)
	SendMessage(context.Context, string, string) error
	ReceiveMessages(context.Context, string, int32) ([]SQSMessage, error)
	DeleteMessage(context.Context, string, string) error
}

type SQSMessage struct {
	MessageID     string
	ReceiptHandle string
	Body          string
}
