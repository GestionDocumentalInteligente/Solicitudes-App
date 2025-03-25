package file

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/teamcubation/sg-file-manager-api/pkg/log"
)

type SQSClient struct {
	sqsService *sqs.SQS
	queueURL   *string
}

func NewSQSClient(region, queueURL string) MessageQueue {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	return &SQSClient{
		queueURL:   aws.String(queueURL),
		sqsService: sqs.New(sess),
	}
}

func (c *SQSClient) SendMessage(ctx context.Context, code string, typeID DocumentTypeID) error {
	logger := log.FromContext(ctx)
	logger.Info("sending sqs message...")

	messageBody := `{"code":"` + code + `","type":"` + fmt.Sprintf("%d", typeID) + `"}`

	_, err := c.sqsService.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    c.queueURL,
		MessageBody: aws.String(messageBody),
	})
	if err != nil {
		logger.Error("error sending sqs message: " + err.Error())
		return err
	}

	return nil
}
