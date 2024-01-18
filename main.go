package main

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/tamirsinai/onboarding-golang/pkg/input"
	"github.com/tamirsinai/onboarding-golang/pkg/output"
	"github.com/tamirsinai/onboarding-golang/pkg/repo"
	"github.com/tamirsinai/onboarding-golang/pkg/scan"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Printf("Error loading AWS config:", err)
		return
	}

	// Create an SQS client
	client := sqs.NewFromConfig(cfg)

	// Specify the URL of the SQS queue
	queueURL := "https://sqs.us-east-1.amazonaws.com/586929748635/onboarding-repo-scan"

	// Receive messages from the queue
	receiveMessages(client, queueURL)

	input, err := input.ReadInputFile()
	if err != nil {
		logger.Error("Error read input file:", zap.Error(err))
		return
	}

	if err := repo.CloneRepositoryToScan(input.CloneUrl); err != nil {
		logger.Error("Error clone repo:", zap.Error(err))
		return
	}

	scan, err := scan.ScanRepoFiles(repo.ClonedProjectsDir, input.Size)
	if err != nil {
		logger.Error("Error scanning repo files:", zap.Error(err))
		return
	}

	if err := output.WriteOutputFile(scan); err != nil {
		logger.Error("Error write output file:", zap.Error(err))
		return
	}

	path, err := filepath.Abs(output.OutputFileName)
	if err != nil {
		logger.Error("Error getting absolute path:", zap.Error(err))
		return
	}
	logger.Info(path)

	if err := repo.DeleteClonedProjectsDir(); err != nil {
		logger.Error("Error with delete repo:", zap.Error(err))
		return
	}
}

func receiveMessages(client *sqs.Client, queueURL string) {
	// Configure the ReceiveMessageInput
	input := &sqs.ReceiveMessageInput{
		QueueUrl:              &queueURL,
		MaxNumberOfMessages:   1,  // Maximum number of messages to receive
		WaitTimeSeconds:       20, // Long polling timeout (adjust as needed)
		MessageAttributeNames: []string{"All"},
	}

	// Receive messages from the queue
	for {
		resp, err := client.ReceiveMessage(context.TODO(), input)
		if err != nil {
			fmt.Println("Error receiving messages:", err)
			return
		}

		// Process received messages
		for _, message := range resp.Messages {
			fmt.Println("Message ID:", aws.ToString(message.MessageId))
			fmt.Println("Receipt Handle:", aws.ToString(message.ReceiptHandle))
			fmt.Println("Body:", aws.ToString(message.Body))

			// Process the message as needed

			// Delete the message from the queue after processing
			deleteMessage(client, queueURL, *message.ReceiptHandle)
		}

		// Sleep for a while before checking for new messages again
		time.Sleep(5 * time.Second)
	}
}

func deleteMessage(client *sqs.Client, queueURL, receiptHandle string) {
	// Configure the DeleteMessageInput
	input := &sqs.DeleteMessageInput{
		QueueUrl:      &queueURL,
		ReceiptHandle: &receiptHandle,
	}

	// Delete the message from the queue
	_, err := client.DeleteMessage(context.TODO(), input)
	if err != nil {
		fmt.Println("Error deleting message:", err)
		return
	}

	fmt.Println("Message deleted.")
}
