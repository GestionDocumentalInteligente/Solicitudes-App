package sdkaws

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/teamcubation/sg-backend/pkg/aws/localstack/defs"
)

type sqsClient struct {
	config    defs.Config
	sqsClient *sqs.Client
}

func (s *stack) NewSQSClient() defs.SqsClient {
	return &sqsClient{
		config: s.config,
		sqsClient: sqs.NewFromConfig(s.awsConfig, func(o *sqs.Options) {
			o.BaseEndpoint = aws.String(s.config.GetLocalStackEndpoint())
		}),
	}
}

// GetOrCreateQueueURL asegura que la cola exista y devuelve su URL
func (q *sqsClient) GetOrCreateQueueURL(ctx context.Context, queueName string) (string, error) {
	// Intentar obtener la URL de la cola
	out, err := q.sqsClient.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err == nil {
		return *out.QueueUrl, nil
	}

	// Si la cola no existe, crearla
	log.Printf("Queue %s not found. Creating it now...", queueName)
	createOut, err := q.sqsClient.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create queue %s: %w", queueName, err)
	}

	return *createOut.QueueUrl, nil
}

// SendMessage sends a message to the specified SQS queue URL
func (q *sqsClient) SendMessage(ctx context.Context, queueURL string, messageBody string) error {
	// Prepare the message input
	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(messageBody),
	}

	// Send the message
	_, err := q.sqsClient.SendMessage(ctx, input)
	if err != nil {
		return fmt.Errorf("error sending message to SQS: %w", err)
	}

	return nil
}

// ReceiveMessages recibe mensajes de la cola especificada
func (q *sqsClient) ReceiveMessages(ctx context.Context, queueURL string, maxMessages int32) ([]defs.SQSMessage, error) {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: maxMessages,
		WaitTimeSeconds:     10, // Long polling
	}

	resp, err := q.sqsClient.ReceiveMessage(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error receiving messages from SQS: %w", err)
	}

	var messages []defs.SQSMessage
	for _, msg := range resp.Messages {
		messages = append(messages, defs.SQSMessage{
			MessageID:     aws.ToString(msg.MessageId),
			ReceiptHandle: aws.ToString(msg.ReceiptHandle),
			Body:          aws.ToString(msg.Body),
		})
	}

	return messages, nil
}

// DeleteMessage elimina un mensaje de la cola especificada usando su ReceiptHandle
func (q *sqsClient) DeleteMessage(ctx context.Context, queueURL string, receiptHandle string) error {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	}

	_, err := q.sqsClient.DeleteMessage(ctx, input)
	if err != nil {
		return fmt.Errorf("error deleting message from SQS: %w", err)
	}

	return nil
}
